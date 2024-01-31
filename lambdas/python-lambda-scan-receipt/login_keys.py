import xml.etree.ElementTree as ET
import requests
import logging

logging.basicConfig(level=logging.INFO)

def login(cms_file_name_path):
    url = "https://wsaa.afip.gov.ar/ws/services/LoginCms"

    # Leer sincrónicamente el contenido del archivo y omitir la primera y última líneas
    with open(f"/tmp/{cms_file_name_path}", "r") as file:
        lines = file.readlines()
        xml_content = "".join(lines[1:-1])

    headers = {
        "Content-Type": "text/xml;charset=utf-8",
        "SOAPAction": "",
    }

    # Construir la solicitud SOAP con el contenido del archivo
    soap_request = f'''<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsaa="http://wsaa.view.sua.dvadac.desein.afip.gov">
        <soapenv:Header/>
        <soapenv:Body>
            <wsaa:loginCms>
                <wsaa:in0>{xml_content}</wsaa:in0>
            </wsaa:loginCms>
        </soapenv:Body>
    </soapenv:Envelope>'''

    timeout = 30.0  # Timeout total de 30 segundos
    try:
        response = requests.post(url, data=soap_request, headers=headers, timeout=timeout)
    except requests.exceptions.ReadTimeout:
        logging.error("Timeout error when trying to POST to the URL: %s", url)
        raise

    root = ET.fromstring(response.content.decode("utf-8"))

    login_cms_return = root.find(".//{http://wsaa.view.sua.dvadac.desein.afip.gov}loginCmsReturn")

    if login_cms_return is None or login_cms_return.text is None:
        raise ValueError("Failed to find loginCmsReturn in the response")

    inner_root = ET.fromstring(login_cms_return.text)
    credentials = inner_root.find(".//credentials")

    if credentials is None:
        return None, None

    token = credentials.find("token").text
    sign = credentials.find("sign").text

    return token, sign
