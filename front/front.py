import json
import os

import google.auth.transport.requests
import google.oauth2.id_token
import pandas as pd
import requests
import streamlit as st


def run() -> None:
    with st.form("gourmet-search"):
        keyword = st.text_input("keyword", value="新宿")
        submitted = st.form_submit_button("送信")
        if submitted:
            st.write("keyword")
            shop_locations = gourmet_search(keyword)
            st.write(shop_locations)
            df = pd.DataFrame.from_dict(shop_locations)
            st.dataframe(df)
            st.map(df)


def gourmet_search(keyword: str) -> list[dict]:
    back_url = os.environ["BACK_URL"]
    id_token = issue_id_token(back_url)
    headers = {"Authorization": f"Bearer {id_token}"}
    params = {"keyword": keyword}
    resp = requests.get(back_url, headers=headers, params=params)
    if resp.status_code != 200:
        msg = {"status_code": resp.status_code, "text": resp.text}
        st.error(json.dumps(msg))
        return list()
    return resp.json()


def issue_id_token(audience: str) -> str:
    auth_req = google.auth.transport.requests.Request()
    id_token = google.oauth2.id_token.fetch_id_token(auth_req, audience)
    return id_token


run()
