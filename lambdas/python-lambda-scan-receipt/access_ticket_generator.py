import aiofiles
from datetime import datetime, timedelta


async def generate_access_ticket(access_ticket_name):
    uniqueId = 190926
    argentina_offset = timedelta(hours=-3)
    now = datetime.utcnow() + argentina_offset

    # Calculate the expiration time (1 day later)
    expiration_time = now + timedelta(hours=22)

    # Format the dates in the desired format
    generation_time_str = now.strftime('%Y-%m-%dT%H:%M:%S')
    expiration_time_str = expiration_time.strftime('%Y-%m-%dT%H:%M:%S')

    # Create the XML content for the ticket
    xml_content = f'''<loginTicketRequest>
        <header>
            <uniqueId>{uniqueId}</uniqueId>
            <generationTime>{generation_time_str}</generationTime>
            <expirationTime>{expiration_time_str}</expirationTime>
        </header>
        <service>ws_sr_constancia_inscripcion</service>
    </loginTicketRequest>'''

    # Save the XML content to a file asynchronously
    async with aiofiles.open(f'/tmp/{access_ticket_name}', 'w') as file:
        await file.write(xml_content)

    return datetime.utcnow()
