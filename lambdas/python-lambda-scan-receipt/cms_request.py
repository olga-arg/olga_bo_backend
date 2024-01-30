import asyncio
import os


async def get_cms_request(access_ticket_name, cms_file_name_path, prod_certificate_path, private_key_path):
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
        process = await asyncio.create_subprocess_exec(*command, stdout=asyncio.subprocess.PIPE, stderr=asyncio.subprocess.PIPE)
        stdout, stderr = await process.communicate()

        if process.returncode == 0:
            return True  # Return a success indicator
        else:
            print(f"Subprocess error: {stderr.decode()}")  # Optionally print the error for debugging
            return False  # Return a failure indicator
    except Exception as e:
        print(f"Error occurred: {e}")  # Optionally print the exception for debugging
        return False  # Return a failure indicator
