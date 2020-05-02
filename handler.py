from typing import Dict
from json import dumps as JSONResponse


def handler(request_args: Dict[str, str]) -> str:
    print(request_args)
    return JSONResponse(request_args)

