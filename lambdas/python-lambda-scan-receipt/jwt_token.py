import os
import requests
from jose import jwt
import boto3
from botocore.exceptions import ClientError

def extract_email_and_company_id_from_token(request):
    auth_header = request['headers'].get('authorization', None)
    print("authHeader: ", auth_header)

    if not auth_header or ' ' not in auth_header:
        return "", "", True

    _, token = auth_header.split(' ', 1)

    user_pool_id = os.getenv('USER_POOL_ID')
    print("upi: ", user_pool_id)

    pub_key_url = f"https://cognito-idp.us-east-1.amazonaws.com/{user_pool_id}/.well-known/jwks.json"
    print("pb: ", pub_key_url)

    response = requests.get(pub_key_url)
    key_set = response.json()

    try:
        decoded_token = jwt.decode(token, key_set)
        print("token: ", decoded_token)
    except jwt.JWTError as e:
        print("JWT error: ", str(e))
        return "", "", True

    username = decoded_token.get('username')

    try:
        cognito_client = boto3.client('cognito-idp', region_name="us-east-1")
        user_data = cognito_client.admin_get_user(
            UserPoolId=user_pool_id,
            Username=username
        )
    except ClientError as e:
        print("Cognito client error: ", e)
        return "", "", True

    email = next((attr['Value'] for attr in user_data['UserAttributes'] if attr['Name'] == 'email'), "")
    company_id = next((attr['Value'] for attr in user_data['UserAttributes'] if attr['Name'] == 'name'), "")

    return email, company_id, False
