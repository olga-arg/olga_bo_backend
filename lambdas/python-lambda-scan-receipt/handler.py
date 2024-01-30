from .scanner import get_receipt_info
from .request_proof_of_registration import request_registration, afip_categories
from .afip_request_credentials import AfipCredentials

import boto3
from botocore.exceptions import ClientError
import asyncio
import json
import os


async def async_lambda_handler(event, context):
    ssm_client = boto3.client('ssm')
    s3 = boto3.client('s3', region_name='us-east-1',
                      aws_access_key_id=os.environ['S3_USER_AC'],
                      aws_secret_access_key=os.environ['S3_USER_SAC'])
    olga_cuit = ssm_client.get_parameter(Name='olga_cuit')['Parameter']['Value']
    mindee_api_key = ssm_client.get_parameter(Name='mindee_api_key')['Parameter']['Value']
    receipt_key = json.loads(event.get('body')).get('receipt_key')
    s3bucket_name = ssm_client.get_parameter(Name='s3_bucket_name')['Parameter']['Value']

    cuit = json.loads(event.get('body')).get('cuit')
    if cuit:
        try:
            afip_creds = await AfipCredentials.create()
            activity_id, company_name = await request_registration(afip_creds.TOKEN, afip_creds.SIGN, olga_cuit, cuit)

        except ValueError as e:
            # If there's an error with request_registration, set default values and log the error
            activity_id, company_name = None, None
            print(f"Error with request_registration: {str(e)}")  # This logs the error for debugging

        if company_name:
            if "SOCIEDAD ANONIMA INDUSTRIAL COMERCIAL FINANCIERA INMOBILIARIA" in company_name:
                company_name = company_name.split("SOCIEDAD ANONIMA INDUSTRIAL COMERCIAL FINANCIERA INMOBILIARIA")[0]
                company_name += 'S.A.I.C.F. E I.'
            elif "SOCIEDAD ANONIMA" in company_name:
                company_name = company_name.split("SOCIEDAD ANONIMA")[0]
                company_name += 'S.A.'
            elif "SOCIEDAD RESPONSABILIDAD LIMITADA" in company_name:
                company_name = company_name.split("SOCIEDAD RESPONSABILIDAD LIMITADA")[0]
                company_name += 'S.R.L.'
            elif "SOCIEDAD DE RESPONSABILIDAD LIMITADA" in company_name:
                company_name = company_name.split("SOCIEDAD DE RESPONSABILIDAD LIMITADA")[0]
                company_name += 'S.R.L.'
            return {
                "statusCode": 200,
                "headers": {
                    "Content-Type": "application/json"
                },
                "body": company_name
            }
        else:
            return {
                "statusCode": 404,
                "headers": {
                    "Content-Type": "application/json"
                },
                "body": "CUIT Not Found"
            }

    try:
        file_metadata = s3.head_object(Bucket=s3bucket_name, Key=receipt_key)
    except ClientError as e:
        error_code = int(e.response['Error']['Code'])
        if error_code == 404 or error_code == 403: # 404 means not found, 403 means forbidden (which can also imply not found in S3)
            return {
                "statusCode": 404,  # 404 is the standard code for not found
                "headers": {
                    "Content-Type": "application/json"
                },
                "body": json.dumps({"error": "Receipt key not found in the S3 bucket"})
            }
        else:
            # Some other unexpected S3 error
            return {
                "statusCode": 500,
                "headers": {
                    "Content-Type": "application/json"
                },
                "body": json.dumps({"error": f"S3 error: {e.response['Error']['Message']}"})
            }

    # Check if file size is greater than 10MB
    file_size = file_metadata['ContentLength']  # size in bytes
    if file_size > 10 * 1024 * 1024:  # 10MB in bytes
        return {
            "statusCode": 413,  # 413 is the standard code for Payload Too Large
            "headers": {
                "Content-Type": "application/json"
            },
            "body": json.dumps({"error": "file too big"})
        }
    file_url = s3.generate_presigned_url('get_object', Params={'Bucket': s3bucket_name, 'Key': receipt_key}, ExpiresIn=60)

    try:
        receipt_info = await get_receipt_info(mindee_api_key, file_url)

    except ValueError as e:
        return {
            "statusCode": 400,  # 400 is the standard code for bad requests
            "headers": {
                "Content-Type": "application/json"
            },
            "body": json.dumps({"error": str(e)})
        }

    try:
        afip_creds = await AfipCredentials.create()
        activity_id, company_name = await request_registration(afip_creds.TOKEN, afip_creds.SIGN, olga_cuit, receipt_info['cuit_number'])

    except ValueError as e:
        # If there's an error with request_registration, set default values and log the error
        activity_id, company_name = None, None
        print(f"Error with request_registration: {str(e)}")  # This logs the error for debugging

    if company_name:
        receipt_info['business_name'] = company_name

    if not activity_id:
        receipt_info['category'] = 'Otros'
    else:
        category = afip_categories(activity_id)
        receipt_info['category'] = category

    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps(receipt_info)
    }


def main(event, context):
    return asyncio.run(async_lambda_handler(event, context))
