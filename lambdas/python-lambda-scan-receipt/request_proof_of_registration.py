import xml.etree.ElementTree as ET
import requests
import time

def request_registration(token, sign, olga_cuit, company_cuit):
    url = "https://aws.afip.gov.ar/sr-padron/webservices/personaServiceA5"
    headers = {
        "Content-Type": "text/xml;charset=utf-8",
        "SOAPAction": "",
    }
    soap_request = f"""<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:a5='http://a5.soap.ws.server.puc.sr/'>
       <soapenv:Header/>
       <soapenv:Body>
          <a5:getPersona>
             <token>{token}</token>
             <sign>{sign}</sign>
             <cuitRepresentada>{olga_cuit}</cuitRepresentada>
             <idPersona>{company_cuit}</idPersona>
          </a5:getPersona>
       </soapenv:Body>
    </soapenv:Envelope>"""

    max_retries = 3
    backoff_factor = 1.5
    timeout = 30.0

    for attempt in range(max_retries):
        try:
            response = requests.post(url, data=soap_request, headers=headers, timeout=timeout)
            root = ET.fromstring(response.content)

            # Handle error message
            error_message = root.find(".//faultstring")
            if error_message is not None:
                if "No existe persona con ese Id" in error_message.text:
                    print(
                        "Afip response: No existe persona con ese Id: ", company_cuit)
                else:
                    print("Afip response: ", error_message.text)
                return None, None

            # Check if activities exist
            activities = root.findall(".//actividad")
            if activities:
                for activity in activities:
                    if activity.find('orden').text == "1":
                        activity_id = activity.find('idActividad').text
                        break
            else:
                activity_id = None

            person_element = root.findall(".//tipoPersona")
            if person_element:
                person = person_element[0].text
            else:
                person = None

            if person == "FISICA":
                comapny_name_element = root.find(".//nombre")
                company_last_name_element = root.find(".//apellido")
                company_name = f"{comapny_name_element.text} {company_last_name_element.text}" if comapny_name_element is not None else None
            else:
                company_name_element = root.find(".//razonSocial")
                company_name = company_name_element.text if company_name_element is not None else None

            return activity_id, company_name

        except requests.exceptions.ReadTimeout:
            if attempt < max_retries - 1:
                sleep_duration = backoff_factor * (2 ** attempt)
                time.sleep(sleep_duration)
            else:
                return None, None
    return None, None

def afip_categories(id):
    categories = {'011111': 'Comidas y Bebidas', '011112': 'Comidas y Bebidas', '011119': 'Comidas y Bebidas', '011121': 'Comidas y Bebidas', '011129': 'Comidas y Bebidas', '011130': 'Comidas y Bebidas', '011211': 'Comidas y Bebidas', '011291': 'Comidas y Bebidas', '011299': 'Comidas y Bebidas', '011310': 'Comidas y Bebidas', '011321': 'Comidas y Bebidas', '011329': 'Comidas y Bebidas', '011331': 'Comidas y Bebidas', '011341': 'Comidas y Bebidas', '011342': 'Comidas y Bebidas', '011400': 'Comidas y Bebidas', '011501': 'Comidas y Bebidas', '011509': 'Comidas y Bebidas', '011911': 'Comidas y Bebidas', '011912': 'Comidas y Bebidas', '011990': 'Comidas y Bebidas', '012110': 'Comidas y Bebidas', '012121': 'Comidas y Bebidas', '012200': 'Comidas y Bebidas', '012311': 'Comidas y Bebidas', '012319': 'Comidas y Bebidas', '012320': 'Comidas y Bebidas', '012410': 'Comidas y Bebidas', '012420': 'Comidas y Bebidas', '012490': 'Comidas y Bebidas', '012510': 'Comidas y Bebidas', '012590': 'Comidas y Bebidas', '012600': 'Comidas y Bebidas', '012701': 'Comidas y Bebidas', '012709': 'Comidas y Bebidas', '012800': 'Comidas y Bebidas', '012900': 'Comidas y Bebidas', '013011': 'Comidas y Bebidas', '013012': 'Comidas y Bebidas', '013013': 'Comidas y Bebidas', '013019': 'Comidas y Bebidas', '013020': 'Indumentaria', '014113': 'Comidas y Bebidas', '014114': 'Comidas y Bebidas', '014115': 'Comisiones y Cargos', '014121': 'Comidas y Bebidas', '014211': 'Otros', '014300': 'Otros', '014410': 'Comidas y Bebidas', '014420': 'Comidas y Bebidas', '014430': 'Comidas y Bebidas', '014440': 'Otros', '014510': 'Otros', '014520': 'Otros', '014610': 'Comidas y Bebidas', '014620': 'Comidas y Bebidas', '014710': 'Otros', '014720': 'Otros', '014810': 'Otros', '014820': 'Otros', '014910': 'Otros', '014920': 'Otros', '014930': 'Otros', '014990': 'Shopping', '016111': 'Cuentas y Servicios', '016112': 'Transporte', '016113': 'Cuentas y Servicios', '016119': 'Cuentas y Servicios', '016120': 'Cuentas y Servicios', '016130': 'Cuentas y Servicios', '016140': 'Cuentas y Servicios', '016150': 'Cuentas y Servicios', '016190': 'Cuentas y Servicios', '016210': 'Cuentas y Servicios', '016220': 'Cuentas y Servicios', '016230': 'Cuentas y Servicios', '016291': 'Cuentas y Servicios', '016292': 'Otros', '016299': 'Cuentas y Servicios', '017010': 'Otros', '017020': 'Cuentas y Servicios', '021010': 'Otros', '021020': 'Otros', '021030': 'Otros', '022010': 'Shopping', '022020': 'Shopping', '024010': 'Cuentas y Servicios', '024020': 'Cuentas y Servicios', '031110': 'Otros', '031120': 'Shopping', '031130': 'Otros', '031200': 'Otros', '031300': 'Cuentas y Servicios', '032000': 'Otros', '051000': 'Otros', '052000': 'Otros', '061000': 'Otros', '062000': 'Otros', '071000': 'Otros', '072100': 'Otros', '072910': 'Otros', '072990': 'Otros', '081100': 'Otros', '081200': 'Otros', '081300': 'Otros', '081400': 'Otros', '089110': 'Inversiones', '089120': 'Shopping', '089200': 'Otros', '089300': 'Otros', '089900': 'Otros', '091000': 'Cuentas y Servicios', '099000': 'Cuentas y Servicios', '101011': 'Comidas y Bebidas', '101012': 'Comidas y Bebidas', '101013': 'Comidas y Bebidas', '101020': 'Comidas y Bebidas', '101030': 'Otros', '101040': 'Comidas y Bebidas', '101091': 'Otros', '101099': 'Comidas y Bebidas', '102001': 'Comidas y Bebidas', '102002': 'Comidas y Bebidas', '102003': 'Comidas y Bebidas', '103011': 'Comidas y Bebidas', '103012': 'Otros', '103020': 'Comidas y Bebidas', '103030': 'Comidas y Bebidas', '103091': 'Hogar', '103099': 'Comidas y Bebidas', '104011': 'Otros', '104012': 'Otros', '104013': 'Otros', '104020': 'Otros', '105010': 'Comidas y Bebidas', '105020': 'Comidas y Bebidas', '105030': 'Otros', '105090': 'Shopping', '106110': 'Comidas y Bebidas', '106120': 'Comidas y Bebidas', '106131': 'Comidas y Bebidas', '106139': 'Comidas y Bebidas', '106200': 'Comidas y Bebidas', '107110': 'Otros', '107121': 'Shopping', '107129': 'Shopping', '107200': 'Otros', '107301': 'Otros', '107309': 'Shopping', '107410': 'Otros', '107420': 'Otros', '107500': 'Comidas y Bebidas', '107911': 'Comidas y Bebidas', '107912': 'Otros', '107920': 'Hogar', '107930': 'Otros', '107991': 'Otros', '107992': 'Otros', '107999': 'Shopping', '108000': 'Comidas y Bebidas', '109000': 'Comidas y Bebidas', '110100': 'Comidas y Bebidas', '110211': 'Otros', '110212': 'Comidas y Bebidas', '110290': 'Comidas y Bebidas', '110300': 'Comidas y Bebidas', '110411': 'Otros', '110412': 'Otros', '110420': 'Comidas y Bebidas', '110491': 'Otros', '110492': 'Comidas y Bebidas', '120010': 'Hogar', '120091': 'Otros', '120099': 'Shopping', '131110': 'Hogar', '131120': 'Hogar', '131131': 'Otros', '131132': 'Otros', '131139': 'Otros', '131201': 'Otros', '131202': 'Otros', '131209': 'Otros', '131300': 'Shopping', '139100': 'Otros', '139201': 'Otros', '139202': 'Indumentaria', '139203': 'Shopping', '139204': 'Shopping', '139209': 'Shopping', '139300': 'Otros', '139400': 'Otros', '139900': 'Shopping', '141110': 'Indumentaria', '141120': 'Indumentaria', '141130': 'Otros', '141140': 'Otros', '141191': 'Indumentaria', '141199': 'Otros', '141201': 'Indumentaria', '141202': 'Otros', '142000': 'Shopping', '143010': 'Otros', '143020': 'Shopping', '149000': 'Cuentas y Servicios', '151100': 'Otros', '151200': 'Shopping', '152011': 'Indumentaria', '152021': 'Indumentaria', '152031': 'Indumentaria', '152040': 'Indumentaria', '161001': 'Otros', '161002': 'Otros', '162100': 'Otros', '162201': 'Hogar', '162202': 'Otros', '162300': 'Otros', '162901': 'Otros', '162902': 'Shopping', '162903': 'Shopping', '162909': 'Shopping', '170101': 'Otros', '170102': 'Otros', '170201': 'Otros', '170202': 'Otros', '170910': 'Shopping', '170990': 'Shopping', '181101': 'Suscripciones', '181109': 'Suscripciones', '181200': 'Cuentas y Servicios', '182000': 'Otros', '191000': 'Shopping', '192000': 'Shopping', '201110': 'Salud y cuidado personal', '201120': 'Otros', '201130': 'Otros', '201140': 'Transporte', '201180': 'Otros', '201190': 'Otros', '201210': 'Otros', '201220': 'Transporte', '201300': 'Inversiones', '201401': 'Otros', '201409': 'Otros', '202101': 'Shopping', '202200': 'Shopping', '202311': 'Hogar', '202312': 'Otros', '202320': 'Salud y cuidado personal', '202906': 'Shopping', '202907': 'Otros', '202908': 'Shopping', '203000': 'Otros', '204000': 'Cuentas y Servicios', '210010': 'Shopping', '210020': 'Mascotas', '210030': 'Otros', '210090': 'Shopping', '221110': 'Otros', '221120': 'Otros', '221901': 'Otros', '221909': 'Shopping', '222010': 'Otros', '222090': 'Shopping', '231010': 'Otros', '231020': 'Otros', '231090': 'Shopping', '239100': 'Shopping', '239201': 'Otros', '239202': 'Otros', '239209': 'Shopping', '239310': 'Shopping', '239391': 'Otros', '239399': 'Shopping', '239410': 'Otros', '239421': 'Otros', '239422': 'Otros', '239510': 'Otros', '239591': 'Otros', '239592': 'Hogar', '239593': 'Shopping', '239600': 'Otros', '239900': 'Shopping', '241001': 'Otros', '241009': 'Shopping', '242010': 'Otros', '242090': 'Shopping', '243100': 'Otros', '243200': 'Otros', '251101': 'Otros', '251102': 'Shopping', '251200': 'Otros', '251300': 'Otros', '252000': 'Otros', '259100': 'Otros', '259200': 'Otros', '259301': 'Indumentaria', '259302': 'Shopping', '259309': 'Shopping', '259910': 'Otros', '259991': 'Otros', '259992': 'Otros', '259993': 'Shopping', '259999': 'Shopping', '261000': 'Electrónica', '262000': 'Shopping', '263000': 'Otros', '264000': 'Shopping', '265101': 'Otros', '265102': 'Otros', '265200': 'Otros', '266010': 'Electrónica', '266090': 'Otros', '267001': 'Indumentaria', '267002': 'Indumentaria', '268000': 'Otros', '271010': 'Otros', '271020': 'Otros', '272000': 'Otros', '273110': 'Salud y cuidado personal', '273190': 'Otros', '274000': 'Otros', '275010': 'Otros', '275020': 'Indumentaria', '275091': 'Otros', '275092': 'Otros', '275099': 'Otros', '279000': 'Otros', '281100': 'Transporte', '281201': 'Otros', '281301': 'Otros', '281400': 'Otros', '281500': 'Hogar', '281600': 'Otros', '281700': 'Otros', '281900': 'Otros', '282110': 'Otros', '282120': 'Otros', '282130': 'Otros', '282200': 'Otros', '282300': 'Otros', '282400': 'Hogar', '282500': 'Comidas y Bebidas', '282600': 'Shopping', '282901': 'Otros', '282909': 'Otros', '291000': 'Transporte', '292000': 'Transporte', '293011': 'Otros', '293090': 'Indumentaria', '301100': 'Hogar', '301200': 'Hogar', '302000': 'Transporte', '303000': 'Hogar', '309100': 'Transporte', '309200': 'Transporte', '309900': 'Transporte', '310010': 'Hogar', '310020': 'Hogar', '310030': 'Otros', '321011': 'Shopping', '321012': 'Otros', '321020': 'Otros', '322001': 'Entretenimiento', '323001': 'Entretenimiento', '324000': 'Entretenimiento', '329010': 'Shopping', '329020': 'Otros', '329030': 'Otros', '329040': 'Indumentaria', '329090': 'Otros', '331101': 'Hogar', '331210': 'Hogar', '331220': 'Hogar', '331290': 'Hogar', '331400': 'Hogar', '331900': 'Hogar', '332000': 'Otros', '351110': 'Otros', '351120': 'Otros', '351130': 'Otros', '351190': 'Otros', '351201': 'Transporte', '351310': 'Shopping', '351320': 'Otros', '352010': 'Otros', '352020': 'Transporte', '353001': 'Otros', '360010': 'Otros', '360020': 'Otros', '370000': 'Cuentas y Servicios', '381100': 'Transporte', '381200': 'Transporte', '382010': 'Otros', '382020': 'Otros', '390000': 'Cuentas y Servicios', '410011': 'Hogar', '410021': 'Hogar', '421000': 'Hogar', '422100': 'Otros', '422200': 'Hogar', '429010': 'Hogar', '429090': 'Hogar', '431100': 'Otros', '431210': 'Hogar', '432110': 'Transporte', '432190': 'Otros', '432200': 'Otros', '432910': 'Otros', '432920': 'Otros', '432990': 'Otros', '433010': 'Otros', '433020': 'Otros', '433030': 'Otros', '433040': 'Hogar', '433090': 'Otros', '439100': 'Hogar', '439910': 'Otros', '439990': 'Hogar', '451110': 'Shopping', '451190': 'Transporte', '451210': 'Shopping', '451290': 'Transporte', '452101': 'Transporte', '452210': 'Hogar', '452220': 'Hogar', '452300': 'Hogar', '452401': 'Hogar', '452500': 'Otros', '452600': 'Hogar', '452700': 'Hogar', '452800': 'Hogar', '452910': 'Hogar', '452990': 'Hogar', '453100': 'Indumentaria', '453210': 'Shopping', '453220': 'Shopping', '453291': 'Indumentaria', '453292': 'Indumentaria', '454010': 'Indumentaria', '454020': 'Transporte', '461011': 'Comidas y Bebidas', '461012': 'Comisiones y Cargos', '461013': 'Comidas y Bebidas', '461014': 'Comidas y Bebidas', '461019': 'Comisiones y Cargos', '461021': 'Comidas y Bebidas', '461022': 'Comidas y Bebidas', '461029': 'Comisiones y Cargos', '461031': 'Comidas y Bebidas', '461032': 'Comidas y Bebidas', '461039': 'Comidas y Bebidas', '461040': 'Transporte', '461092': 'Hogar', '461093': 'Comisiones y Cargos', '461094': 'Comisiones y Cargos', '461095': 'Educación', '461099': 'Comisiones y Cargos', '462110': 'Otros', '462120': 'Shopping', '462131': 'Comidas y Bebidas', '462132': 'Otros', '462190': 'Shopping', '462201': 'Shopping', '462209': 'Shopping', '463111': 'Shopping', '463112': 'Comidas y Bebidas', '463121': 'Comidas y Bebidas', '463129': 'Shopping', '463130': 'Comidas y Bebidas', '463140': 'Comidas y Bebidas', '463151': 'Shopping', '463152': 'Shopping', '463153': 'Shopping', '463154': 'Comidas y Bebidas', '463159': 'Shopping', '463160': 'Shopping', '463170': 'Comidas y Bebidas', '463180': 'Supermercado', '463191': 'Comidas y Bebidas', '463199': 'Shopping', '463211': 'Comidas y Bebidas', '463212': 'Comidas y Bebidas', '463219': 'Comidas y Bebidas', '463220': 'Comidas y Bebidas', '463300': 'Shopping', '464111': 'Shopping', '464112': 'Shopping', '464113': 'Indumentaria', '464114': 'Shopping', '464119': 'Shopping', '464121': 'Shopping', '464122': 'Shopping', '464129': 'Indumentaria', '464130': 'Indumentaria', '464141': 'Shopping', '464142': 'Shopping', '464149': 'Shopping', '464150': 'Indumentaria', '464211': 'Educación', '464212': 'Suscripciones', '464221': 'Shopping', '464222': 'Shopping', '464223': 'Shopping', '464310': 'Shopping', '464320': 'Salud y cuidado personal', '464330': 'Shopping', '464340': 'Mascotas', '464410': 'Salud y cuidado personal',
                  '464420': 'Shopping', '464501': 'Hogar', '464502': 'Shopping', '464610': 'Hogar', '464620': 'Shopping', '464631': 'Shopping', '464632': 'Shopping', '464920': 'Hogar', '464930': 'Shopping', '464940': 'Transporte', '464950': 'Entretenimiento', '464991': 'Shopping', '464999': 'Shopping', '465100': 'Indumentaria', '465210': 'Shopping', '465220': 'Electrónica', '465310': 'Shopping', '465320': 'Comidas y Bebidas', '465330': 'Indumentaria', '465340': 'Shopping', '465350': 'Shopping', '465360': 'Shopping', '465390': 'Shopping', '465400': 'Shopping', '465500': 'Transporte', '465610': 'Hogar', '465690': 'Hogar', '465910': 'Shopping', '465920': 'Shopping', '465930': 'Shopping', '465990': 'Shopping', '466110': 'Transporte', '466121': 'Otros', '466129': 'Transporte', '466200': 'Shopping', '466310': 'Shopping', '466320': 'Hogar', '466330': 'Shopping', '466340': 'Shopping', '466350': 'Shopping', '466360': 'Hogar', '466370': 'Hogar', '466391': 'Hogar', '466399': 'Hogar', '466910': 'Shopping', '466920': 'Shopping', '466931': 'Shopping', '466932': 'Inversiones', '466939': 'Shopping', '466940': 'Shopping', '466990': 'Shopping', '469010': 'Shopping', '469090': 'Shopping', '471110': 'Supermercado', '471120': 'Supermercado', '471130': 'Supermercado', '471190': 'Shopping', '471900': 'Comidas y Bebidas', '472111': 'Shopping', '472112': 'Shopping', '472120': 'Shopping', '472130': 'Comidas y Bebidas', '472140': 'Comidas y Bebidas', '472150': 'Comidas y Bebidas', '472160': 'Comidas y Bebidas', '472171': 'Shopping', '472172': 'Shopping', '472190': 'Shopping', '472200': 'Comidas y Bebidas', '472300': 'Shopping', '473000': 'Transporte', '474010': 'Indumentaria', '474020': 'Shopping', '475110': 'Shopping', '475120': 'Hogar', '475190': 'Shopping', '475210': 'Shopping', '475220': 'Hogar', '475230': 'Shopping', '475240': 'Shopping', '475250': 'Hogar', '475260': 'Shopping', '475270': 'Hogar', '475290': 'Hogar', '475300': 'Hogar', '475410': 'Hogar', '475420': 'Shopping', '475430': 'Shopping', '475440': 'Shopping', '475490': 'Hogar', '476110': 'Educación', '476120': 'Suscripciones', '476130': 'Shopping', '476310': 'Shopping', '476320': 'Shopping', '476400': 'Entretenimiento', '477110': 'Indumentaria', '477120': 'Shopping', '477130': 'Shopping', '477140': 'Shopping', '477150': 'Shopping', '477190': 'Indumentaria', '477210': 'Shopping', '477220': 'Indumentaria', '477230': 'Indumentaria', '477290': 'Shopping', '477310': 'Shopping', '477320': 'Salud y cuidado personal', '477330': 'Shopping', '477410': 'Salud y cuidado personal', '477420': 'Shopping', '477430': 'Shopping', '477440': 'Inversiones', '477450': 'Hogar', '477460': 'Shopping', '477470': 'Mascotas', '477480': 'Shopping', '477490': 'Shopping', '477810': 'Hogar', '477820': 'Educación', '477830': 'Shopping', '477840': 'Shopping', '477890': 'Transporte', '478010': 'Comidas y Bebidas', '478090': 'Electrónica', '479101': 'Shopping', '479109': 'Shopping', '479900': 'Shopping', '491110': 'Transporte', '491120': 'Transporte', '491200': 'Transporte', '492110': 'Transporte', '492120': 'Transporte', '492130': 'Transporte', '492140': 'Transporte', '492150': 'Transporte', '492160': 'Transporte', '492170': 'Transporte', '492180': 'Transporte', '492190': 'Transporte', '492210': 'Cuentas y Servicios', '492221': 'Transporte', '492229': 'Transporte', '492230': 'Transporte', '492240': 'Transporte', '492250': 'Transporte', '492280': 'Transporte', '492290': 'Transporte', '493110': 'Transporte', '493120': 'Transporte', '493200': 'Transporte', '501100': 'Transporte', '501200': 'Transporte', '502101': 'Transporte', '502200': 'Transporte', '511000': 'Transporte', '512000': 'Transporte', '521010': 'Transporte', '521020': 'Cuentas y Servicios', '521030': 'Transporte', '522010': 'Cuentas y Servicios', '522020': 'Cuentas y Servicios', '522091': 'Cuentas y Servicios', '522092': 'Cuentas y Servicios', '522099': 'Cuentas y Servicios', '523011': 'Salud y cuidado personal', '523019': 'Transporte', '523020': 'Transporte', '523031': 'Transporte', '523032': 'Cuentas y Servicios', '523039': 'Cuentas y Servicios', '523090': 'Transporte', '524110': 'Transporte', '524120': 'Cuentas y Servicios', '524130': 'Transporte', '524190': 'Transporte', '524210': 'Transporte', '524220': 'Cuentas y Servicios', '524230': 'Cuentas y Servicios', '524290': 'Transporte', '524310': 'Transporte', '524320': 'Cuentas y Servicios', '524330': 'Cuentas y Servicios', '524390': 'Transporte', '530010': 'Cuentas y Servicios', '530090': 'Cuentas y Servicios', '551010': 'Cuentas y Servicios', '551021': 'Cuentas y Servicios', '551022': 'Comidas y Bebidas', '551023': 'Comidas y Bebidas', '551090': 'Cuentas y Servicios', '552000': 'Cuentas y Servicios', '561011': 'Comidas y Bebidas', '561012': 'Comidas y Bebidas', '561013': 'Comidas y Bebidas', '561014': 'Comidas y Bebidas', '561019': 'Comidas y Bebidas', '561020': 'Comidas y Bebidas', '561030': 'Cuentas y Servicios', '561040': 'Comidas y Bebidas', '562010': 'Comidas y Bebidas', '562091': 'Comidas y Bebidas', '562099': 'Comidas y Bebidas', '581100': 'Educación', '581200': 'Otros', '581300': 'Suscripciones', '581900': 'Otros', '591110': 'Otros', '591120': 'Otros', '591200': 'Otros', '591300': 'Otros', '592000': 'Entretenimiento', '601000': 'Otros', '602100': 'Otros', '602200': 'Suscripciones', '602310': 'Suscripciones', '602320': 'Otros', '602900': 'Cuentas y Servicios', '611010': 'Cuentas y Servicios', '611090': 'Cuentas y Servicios', '612000': 'Electrónica', '613000': 'Cuentas y Servicios', '614010': 'Cuentas y Servicios', '614090': 'Cuentas y Servicios', '619000': 'Cuentas y Servicios', '620100': 'Cuentas y Servicios', '620200': 'Cuentas y Servicios', '620300': 'Cuentas y Servicios', '620900': 'Cuentas y Servicios', '631110': 'Otros', '631120': 'Otros', '631190': 'Otros', '631200': 'Otros', '639100': 'Otros', '639900': 'Cuentas y Servicios', '641100': 'Cuentas y Servicios', '641910': 'Cuentas y Servicios', '641920': 'Cuentas y Servicios', '641930': 'Cuentas y Servicios', '641941': 'Cuentas y Servicios', '641942': 'Hogar', '641943': 'Cuentas y Servicios', '642000': 'Cuentas y Servicios', '643001': 'Cuentas y Servicios', '643009': 'Inversiones', '649100': 'Cuentas y Servicios', '649210': 'Cuentas y Servicios', '649220': 'Cuentas y Servicios', '649290': 'Cuentas y Servicios', '649910': 'Supermercado', '649991': 'Cuentas y Servicios', '649999': 'Cuentas y Servicios', '651110': 'Salud y cuidado personal', '651120': 'Cuentas y Servicios', '651130': 'Salud y cuidado personal', '651210': 'Cuentas y Servicios', '651220': 'Cuentas y Servicios', '651310': 'Otros', '651320': 'Cuentas y Servicios', '652000': 'Otros', '653000': 'Mascotas', '661111': 'Supermercado', '661121': 'Supermercado', '661131': 'Cuentas y Servicios', '661910': 'Cuentas y Servicios', '661920': 'Cuentas y Servicios', '661930': 'Cuentas y Servicios', '661991': 'Cuentas y Servicios', '661992': 'Cuentas y Servicios', '661999': 'Cuentas y Servicios', '662010': 'Cuentas y Servicios', '662020': 'Cuentas y Servicios', '662090': 'Cuentas y Servicios', '663000': 'Cuentas y Servicios', '681010': 'Entretenimiento', '681020': 'Cuentas y Servicios', '681098': 'Cuentas y Servicios', '681099': 'Cuentas y Servicios', '682010': 'Cuentas y Servicios', '682091': 'Cuentas y Servicios', '682099': 'Cuentas y Servicios', '691001': 'Cuentas y Servicios', '691002': 'Cuentas y Servicios', '692000': 'Servicios profesionales', '702010': 'Salud y cuidado personal', '702091': 'Cuentas y Servicios', '702092': 'Cuentas y Servicios', '702099': 'Cuentas y Servicios', '711001': 'Hogar', '711002': 'Cuentas y Servicios', '711003': 'Cuentas y Servicios', '711009': 'Cuentas y Servicios', '712000': 'Otros', '721010': 'Otros', '721020': 'Otros', '721030': 'Otros', '721090': 'Otros', '722010': 'Otros', '722020': 'Otros', '731001': 'Salud y cuidado personal', '731009': 'Cuentas y Servicios', '732000': 'Supermercado', '741000': 'Cuentas y Servicios', '742000': 'Cuentas y Servicios', '749001': 'Cuentas y Servicios', '749002': 'Cuentas y Servicios', '749003': 'Cuentas y Servicios', '749009': 'Otros', '750000': 'Mascotas', '771110': 'Transporte', '771190': 'Transporte', '771210': 'Transporte', '771220': 'Transporte', '771290': 'Transporte', '772010': 'Entretenimiento', '772091': 'Otros', '772099': 'Otros', '773010': 'Otros', '773020': 'Otros', '773030': 'Hogar', '773040': 'Electrónica', '773090': 'Otros', '774000': 'Cuentas y Servicios', '780000': 'Otros', '791100': 'Viajes', '791200': 'Viajes', '791901': 'Viajes', '791909': 'Cuentas y Servicios', '801010': 'Transporte', '801020': 'Cuentas y Servicios', '801090': 'Cuentas y Servicios', '811000': 'Cuentas y Servicios', '812010': 'Hogar', '812020': 'Cuentas y Servicios', '812090': 'Hogar', '813000': 'Salud y cuidado personal', '821100': 'Cuentas y Servicios', '821900': 'Hogar', '822000': 'Cuentas y Servicios', '823000': 'Cuentas y Servicios', '829100': 'Cuentas y Servicios', '829200': 'Cuentas y Servicios', '829900': 'Cuentas y Servicios', '841100': 'Cuentas y Servicios', '841200': 'Mascotas', '841300': 'Cuentas y Servicios', '841900': 'Cuentas y Servicios', '842100': 'Cuentas y Servicios', '842200': 'Cuentas y Servicios', '842300': 'Cuentas y Servicios', '842400': 'Cuentas y Servicios', '842500': 'Cuentas y Servicios', '843000': 'Mascotas', '851010': 'Otros', '851020': 'Hogar', '852100': 'Educación', '852200': 'Educación', '853100': 'Educación', '853201': 'Educación', '853300': 'Otros', '854910': 'Educación', '854920': 'Educación', '854930': 'Educación', '854940': 'Educación', '854950': 'Entretenimiento', '854960': 'Educación', '854990': 'Educación', '855000': 'Educación', '861010': 'Salud y cuidado personal', '861020': 'Salud y cuidado personal', '862110': 'Cuentas y Servicios', '862120': 'Cuentas y Servicios', '862130': 'Salud y cuidado personal', '862200': 'Cuentas y Servicios', '863110': 'Cuentas y Servicios', '863120': 'Cuentas y Servicios', '863190': 'Cuentas y Servicios', '863200': 'Cuentas y Servicios', '863300': 'Cuentas y Servicios', '864000': 'Cuentas y Servicios', '869010': 'Cuentas y Servicios', '869090': 'Salud y cuidado personal', '870100': 'Salud y cuidado personal', '870210': 'Cuentas y Servicios', '870220': 'Cuentas y Servicios', '870910': 'Cuentas y Servicios', '870920': 'Cuentas y Servicios', '870990': 'Cuentas y Servicios', '880000': 'Cuentas y Servicios', '900011': 'Entretenimiento', '900021': 'Otros', '900030': 'Entretenimiento', '900040': 'Cuentas y Servicios', '900091': 'Entretenimiento', '910100': 'Educación', '910200': 'Cuentas y Servicios', '910300': 'Entretenimiento', '910900': 'Cuentas y Servicios', '920001': 'Cuentas y Servicios', '920009': 'Entretenimiento', '931010': 'Cuentas y Servicios', '931020': 'Otros', '931030': 'Entretenimiento', '931041': 'Cuentas y Servicios', '931042': 'Cuentas y Servicios', '931050': 'Cuentas y Servicios', '931090': 'Cuentas y Servicios', '939010': 'Entretenimiento', '939020': 'Entretenimiento', '939030': 'Cuentas y Servicios', '939090': 'Cuentas y Servicios', '941100': 'Cuentas y Servicios', '941200': 'Cuentas y Servicios', '942000': 'Cuentas y Servicios', '949100': 'Cuentas y Servicios', '949200': 'Cuentas y Servicios', '949910': 'Salud y cuidado personal', '949920': 'Cuentas y Servicios', '949930': 'Cuentas y Servicios', '949990': 'Cuentas y Servicios', '951100': 'Hogar', '951200': 'Hogar', '952200': 'Shopping', '952300': 'Hogar', '952910': 'Hogar', '952920': 'Hogar', '952990': 'Hogar', '960101': 'Hogar', '960102': 'Hogar', '960201': 'Salud y cuidado personal', '960202': 'Salud y cuidado personal', '960300': 'Cuentas y Servicios', '960910': 'Salud y cuidado personal', '960990': 'Cuentas y Servicios', '970000': 'Hogar', '990000': 'Cuentas y Servicios', '952100': 'Electrónica', '476200': 'Shopping', '464910': 'Shopping', '461091': 'Indumentaria', '431220': 'Otros', '331301': 'Hogar', '014221': 'Otros', '000007': 'Otros', '000008': 'Otros', '000009': 'Otros', '000010': 'Otros', '000011': 'Otros', '000012': 'Otros', '000013': 'Otros'}

    return categories.get(id)
