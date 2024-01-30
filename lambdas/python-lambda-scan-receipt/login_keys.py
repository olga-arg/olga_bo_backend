import xml.etree.ElementTree as ET
import httpx
import aiofiles
import logging

logging.basicConfig(level=logging.INFO)


async def login(cms_file_name_path):
    url = "https://wsaa.afip.gov.ar/ws/services/LoginCms"

    # Asynchronously read the content of the file and omit the first and last lines
    async with aiofiles.open(f"/tmp/{cms_file_name_path}", "r") as file:
        lines = await file.readlines()
        xml_content = "".join(lines[1:-1])

    headers = {
        "Content-Type": "text/xml;charset=utf-8",
        "SOAPAction": "",
    }

    # Construct the SOAP request with the content from the file
    soap_request = f'''<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsaa="http://wsaa.view.sua.dvadac.desein.afip.gov">
        <soapenv:Header/>
        <soapenv:Body>
            <wsaa:loginCms>
                <wsaa:in0>{xml_content}</wsaa:in0>
            </wsaa:loginCms>
        </soapenv:Body>
    </soapenv:Envelope>'''

    timeout = httpx.Timeout(30.0, read=30.0)  # Ejemplo: timeout total de 30 segundos, timeout de lectura de 30 segundos
    async with httpx.AsyncClient(timeout=timeout) as client:
        try:
            response = await client.post(url, data=soap_request, headers=headers)
        except httpx.ReadTimeout:
            logging.error("Timeout error when trying to POST to the URL: %s", url)
            raise

    root = ET.fromstring(response.content.decode("utf-8"))
    print(response.content.decode("utf-8"))

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
