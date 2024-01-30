from access_ticket_generator import generate_access_ticket
from cms_request import get_cms_request
from login_keys import login

from datetime import datetime
import aioboto3
import asyncio


class AfipCredentials:
    def __init__(self):
        self.TOKEN = None
        self.SIGN = None

        self.credentials_timestamp_creation = None
        self.access_ticket_name = None
        self.cms_file_name_path = None
        self.prod_certificate_path = None
        self.private_key_path = None

        self.session = aioboto3.Session()

    @classmethod
    async def create(cls, force):
        instance = cls()
        await instance.get_ssm_parameters()
        await instance.request_afip_credentials(force)
        return instance

    async def request_afip_credentials(self, force):
        async with self.session.client('ssm') as client:
            current_utc_time = int(datetime.utcnow().timestamp())
            # Si han pasado más de 10 horas desde el último login, ya que expira a las 12 horas.
            print("UTC ahora: ", current_utc_time)
            print("Cred timestamp creation: ",
                  self.credentials_timestamp_creation)
            print("UTC ahora - cred : ", current_utc_time -
                  self.credentials_timestamp_creation)
            if current_utc_time > self.credentials_timestamp_creation + 36000 or force:
                await generate_access_ticket(self.access_ticket_name)
                await get_cms_request(self.access_ticket_name, self.cms_file_name_path, self.prod_certificate_path, self.private_key_path)
                await asyncio.sleep(0.2)
                self.TOKEN, self.SIGN = await login(self.cms_file_name_path)

                await client.put_parameter(Name='afip_token', Value=self.TOKEN, Type='String', Tier='Standard', Overwrite=True)
                await client.put_parameter(Name='afip_sign', Value=self.SIGN, Type='String', Tier='Standard', Overwrite=True)
                await client.put_parameter(Name='afip_credentials_timestamp_creation', Type='String', Tier='Standard', Value=str(current_utc_time), Overwrite=True)

    # Continue with the rest of the operations
    async def get_ssm_parameters(self):
        async with self.session.client('ssm') as client:
            response = await client.get_parameter(Name='afip_credentials_timestamp_creation')
            self.credentials_timestamp_creation = int(
                response['Parameter']['Value'])

            response = await client.get_parameter(Name='afip_access_ticket_path')
            self.access_ticket_name = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_cms_file_path')
            self.cms_file_name_path = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_prod_certificate_path')
            self.prod_certificate_path = response['Parameter']['Value']

            response = await client.get_parameter(Name='afip_private_key_path')
            self.private_key_path = response['Parameter']['Value']

            try:
                response = await client.get_parameter(Name='afip_token')
                self.TOKEN = response['Parameter']['Value']

                response = await client.get_parameter(Name='afip_sign')
                self.SIGN = response['Parameter']['Value']
            except:
                pass
