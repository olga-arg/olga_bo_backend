import os
import subprocess

def get_cms_request(access_ticket_name, cms_file_name_path, prod_certificate_path, private_key_path):
    current_path = os.path.dirname(os.path.abspath(__file__))
    prod_certificate_path = os.path.join(current_path, prod_certificate_path)
    private_key_path = os.path.join(current_path, private_key_path)

    command = [
        "openssl",
        "cms",
        "-sign",
        "-in",
        f'/tmp/{access_ticket_name}',
        "-out",
        f'/tmp/{cms_file_name_path}',
        "-signer",
        prod_certificate_path,
        "-inkey",
        private_key_path,
        "-nodetach",
        "-outform",
        "PEM"
    ]

    try:
        process = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=True)
        return True  # Return a success indicator
    except subprocess.CalledProcessError as e:
        print(f"Subprocess error: {e.stderr.decode()}")  # Optionally print the error for debugging
        return False  # Return a failure indicator
    except Exception as e:
        print(f"Error occurred: {e}")  # Optionally print the exception for debugging
        return False  # Return a failure indicator
