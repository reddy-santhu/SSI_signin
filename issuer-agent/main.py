#!/usr/bin/env python3

import requests
import json
import os

ISSUER_AGENT_URL = os.getenv("ISSUER_AGENT_URL", "http://localhost:8002")
LEDGER_URL = os.getenv("LEDGER_URL", "http://localhost:9000")


def create_schema(name, version, attributes):
    url = f"{ISSUER_AGENT_URL}/schemas"
    
    payload = {
        "schema_name": name,
        "schema_version": version,
        "attributes": attributes
    }
    
    response = requests.post(url, json=payload)
    response.raise_for_status()
    
    return response.json()


def create_credential_definition(schema_id, support_revocation=True):
    url = f"{ISSUER_AGENT_URL}/credential-definitions"
    
    payload = {
        "schema_id": schema_id,
        "support_revocation": support_revocation
    }
    
    response = requests.post(url, json=payload)
    response.raise_for_status()
    
    return response.json()


def issue_credential(connection_id, cred_def_id, attributes):
    url = f"{ISSUER_AGENT_URL}/issue-credential/send"
    
    payload = {
        "connection_id": connection_id,
        "credential_definition_id": cred_def_id,
        "credential_proposal": {
            "@type": "issue-credential/1.0/credential-proposal",
            "attributes": attributes
        }
    }
    
    response = requests.post(url, json=payload)
    response.raise_for_status()
    
    return response.json()


if __name__ == "__main__":
    print("Creating KYC schema...")
    schema = create_schema(
        name="KYC Credential",
        version="1.0",
        attributes=["phone_number", "country", "kyc_verified", "issue_date"]
    )
    print(f"Schema created: {schema.get('schema_id')}")
    
    print("Creating credential definition...")
    cred_def = create_credential_definition(
        schema_id=schema.get("schema_id"),
        support_revocation=True
    )
    print(f"Credential definition created: {cred_def.get('credential_definition_id')}")

