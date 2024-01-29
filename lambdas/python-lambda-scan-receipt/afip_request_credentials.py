from .access_ticket_generator import generate_access_ticket
from .cms_request import get_cms_request
from .login_keys import login

from datetime import datetime
import aioboto3
import asyncio


class AfipCredentials:
    def __init__(self):
        self.TOKEN = None
        self.SIGN = None

        self.credentials_timestamp = None
        self.access_ticket_name = None
        self.cms_file_name_path = None
        self.prod_certificate_path = None
        self.private_key_path = None

        self.session = aioboto3.Session()

    @classmethod
    async def create(cls):
        instance = cls()
        await instance.get_ssm_parameters()
        await instance.request_afip_credentials()
        return instance

    async def request_afip_credentials(self):
        async with self.session.client('ssm') as client:
            current_utc_time = int(datetime.utcnow().timestamp())

            # Si han pasado más de 23 horas desde el último login
            if current_utc_time - self.credentials_timestamp > 82800:
                await generate_access_ticket(self.access_ticket_name)
                await get_cms_request(self.access_ticket_name, self.cms_file_name_path, self.prod_certificate_path, self.private_key_path)
                await asyncio.sleep(0.2)
                self.TOKEN, self.SIGN = await login(self.cms_file_name_path)

                await client.put_parameter(Name='afip_token', Value=self.TOKEN, Overwrite=True)
                await client.put_parameter(Name='afip_sign', Value=self.SIGN, Overwrite=True)
                await client.put_parameter(Name='afip_credentials_timestamp', Value=str(current_utc_time), Overwrite=True)
            else:
                self.TOKEN = (await client.get_parameter(Name='afip_token'))['Parameter']['Value']
                self.SIGN = (await client.get_parameter(Name='afip_sign'))['Parameter']['Value']

    # Continue with the rest of the operations
    async def get_ssm_parameters(self):
        async with self.session.client('ssm') as client:
            response = await client.get_parameter(Name='afip_credentials_timestamp')
            self.credentials_timestamp = int(response['Parameter']['Value'])

            response = await client.get_parameter(Name='afip_access_ticket_path')
            self.access_ticket_name = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_cms_file_path')
            self.cms_file_name_path = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_prod_certificate_path')
            self.prod_certificate_path = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_private_key_path')
            self.private_key_path = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_token')
            self.TOKEN = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_sign')
            self.SIGN = response['Parameter']['Value']
