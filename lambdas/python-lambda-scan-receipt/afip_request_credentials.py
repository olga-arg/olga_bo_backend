from access_ticket_generator import generate_access_ticket
from cms_request import get_cms_request
from login_keys import login

from datetime import datetime
import boto3

class AfipCredentials:
    def __init__(self):
        self.TOKEN = None
        self.SIGN = None

        self.credentials_timestamp_creation = None
        self.access_ticket_name = None
        self.cms_file_name_path = None
        self.prod_certificate_path = None
        self.private_key_path = None

        self.session = boto3.Session()

    @classmethod
    def create(cls, force):
        instance = cls()
        instance.get_ssm_parameters()
        instance.request_afip_credentials(force)
        return instance

    def request_afip_credentials(self, force):
        client = self.session.client('ssm')
        current_utc_time = int(datetime.utcnow().timestamp())
        print("UTC ahora: ", current_utc_time)
        print("Cred timestamp creation: ", self.credentials_timestamp_creation)
        print("UTC ahora - cred : ", current_utc_time - self.credentials_timestamp_creation)

        if current_utc_time > self.credentials_timestamp_creation + 36000 or force:
            generate_access_ticket(self.access_ticket_name)
            get_cms_request(self.access_ticket_name, self.cms_file_name_path, self.prod_certificate_path, self.private_key_path)
            self.TOKEN, self.SIGN = login(self.cms_file_name_path)

            client.put_parameter(Name='afip_token', Value=self.TOKEN, Type='String', Tier='Standard', Overwrite=True)
            client.put_parameter(Name='afip_sign', Value=self.SIGN, Type='String', Tier='Standard', Overwrite=True)
            client.put_parameter(Name='afip_credentials_timestamp_creation', Type='String', Tier='Standard', Value=str(current_utc_time), Overwrite=True)

    def get_ssm_parameters(self):
        client = self.session.client('ssm')
        response = client.get_parameter(Name='afip_credentials_timestamp_creation')
        self.credentials_timestamp_creation = int(response['Parameter']['Value'])

        response = client.get_parameter(Name='afip_access_ticket_path')
        self.access_ticket_name = response['Parameter']['Value']

        response = client.get_parameter(Name='afip_cms_file_path')
        self.cms_file_name_path = response['Parameter']['Value']

        response = client.get_parameter(Name='afip_prod_certificate_path')
        self.prod_certificate_path = response['Parameter']['Value']

        response = client.get_parameter(Name='afip_private_key_path')
        self.private_key_path = response['Parameter']['Value']

        try:
            response = client.get_parameter(Name='afip_token')
            self.TOKEN = response['Parameter']['Value']

            response = client.get_parameter(Name='afip_sign')
            self.SIGN = response['Parameter']['Value']
        except:
            pass
