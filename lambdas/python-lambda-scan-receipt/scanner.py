import httpx
import os
import psycopg2
import requests


def correct_amount_format(amount: str) -> str:
    # Eliminar todas las comas y puntos
    amount = amount.replace(',', '').replace('.', '')

    # Insertar un punto para los decimales
    amount = amount[:-2] + '.' + amount[-2:]

    return amount


def correct_amount_format_v2(amount: str) -> str:
    # Si ya tiene un punto, simplemente garantizar dos decimales
    if '.' in amount:
        integer_part, decimal_part = amount.split('.')
        if len(decimal_part) < 2:
            decimal_part += '0' * (2 - len(decimal_part))
        return f"{integer_part}.{decimal_part}"

    # Si no tiene un punto, seguir el proceso original
    return correct_amount_format(amount)


def get_receipt_info(api_key, file_url, company_id):
    try:
        db_user = os.environ['DB_USER']
        db_password = os.environ['DB_PASSWORD']
        db_name = os.environ['DB_NAME']
        db_host_reader = os.environ['DB_HOST_READER']
        db_port = os.environ['DB_PORT']

        # Establecer la conexión con la base de datos
        conn = psycopg2.connect(
            user=db_user,
            password=db_password,
            dbname=db_name,
            host=db_host_reader,
            port=db_port
        )

        # Crear un cursor para ejecutar consultas
        cur = conn.cursor()

        # Ejecutar la consulta SQL
        sql = f"SELECT cuit FROM companies WHERE id = '{company_id}' LIMIT 1"
        cur.execute(sql)

        # Obtener el resultado
        company_cuit = cur.fetchone()
        if company_cuit:
            company_cuit = company_cuit[0]
        
        print('Company cuit: ', company_cuit)

        # Cerrar el cursor y la conexión
        cur.close()
        conn.close()

        response = requests.post(
            'https://api.mindee.net/v1/products/mindee/argentine-expense-receipt/v1/predict',
            headers={"Authorization": f"Token {api_key}"},
            data={'document': file_url},
            timeout=30.0
        )
        response.raise_for_status()  # Raise an error for HTTP errors
        result = response.json()

        # Extracting the necessary fields from the result
        info = {}

        predictions = result['document']['inference']['prediction']

        fields = ['business_name', 'cuit_number', 'receipt_number', 'receipt_or_ticket_type', 'receipt_datetime', 'total_amount']

        for field in fields:
            if predictions[field]['value']:
                content = predictions[field]['value']
                if field == 'cuit_number':
                    content = ''.join(filter(lambda i: i.isdigit(), content))

                if field == 'total_amount':
                    print(type(content), content)
                    content = correct_amount_format_v2(str(content))
                
                info[field] = content
        
        if info['cuit_number'] == company_cuit:
            raise ValueError("Cuit not visible")

        # If 'cuit' key is not present or its length is not 11, raise an error
        if 'cuit_number' not in info or len(info['cuit_number']) != 11:
            raise ValueError("Cuit not visible")

        if 'CONSUMIDOR' or 'FINAL' in info['receipt_or_ticket_type'].upper():
            info['receipt_or_ticket_type'] = 'B'
        
        # Amount cannot be negative
        if float(info['total_amount']) < 0:
            raise ValueError("Negative amount")

        print('Receipt info: ', info)

        return info

    except httpx.HTTPError:
        raise ValueError("HTTP error occurred while fetching receipt info.")
    except KeyError:
        raise ValueError("Unexpected response format when fetching receipt info.")
    except Exception as e:
        raise ValueError(f"An error occurred while fetching receipt info: {str(e)}")
