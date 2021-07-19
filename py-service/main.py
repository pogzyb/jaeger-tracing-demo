import logging
import os
from typing import List, Any

from fastapi import FastAPI, APIRouter, Request, HTTPException, status
from jaeger_client import Config, Tracer, Span

from config import get_config
from schemas import MatrixOut


# ======================
# Setup and Configs
# ======================

# jaeger tracer client
tracer: Tracer


def init_tracer(service_name: str):
    c = Config(
        config={
            'sampler': {
                'type': 'const',
                'param': 1,
            },
            'local_agent': {
                'reporting_host': os.getenv('JAEGER_AGENT_HOST'),
                'reporting_port': os.getenv('JAEGER_AGENT_PORT'),
            },
            'logging': True,
        },
        service_name=service_name,
        validate=True)
    return c.initialize_tracer()


def factory() -> FastAPI:
    cfg = get_config()
    app = FastAPI(
        debug=cfg.debug,
        title=cfg.app_name,
        description=cfg.description,
        version=cfg.app_version)

    app.include_router(api)

    @app.on_event('startup')
    async def init_app_logger():
        handler = logging.StreamHandler()
        handler.setFormatter(logging.Formatter("%(asctime)s - %(levelname)s - %(message)s"))
        logging.getLogger("uvicorn.access").addHandler(handler)

    @app.on_event('startup')
    async def init_jaeger_tracer():
        global tracer
        tracer = init_tracer(cfg.app_name)

    return app


# ======================
# Route Handlers
# ======================

logger = logging.getLogger(__name__)
api = APIRouter()


# def create_matrix(r: int, c: int, span: Span) -> List[Any]:  # no need to pass span!
def create_matrix(r: int, c: int) -> List[Any]:
    with tracer.start_active_span("create-matrix") as scope:
        elements = 0
        outer = []
        for i in range(r):
            inner = []
            for j in range(c):
                inner.append(j)
                elements += 1
            outer.append(inner)
        scope.span.log_kv({"total-elements": elements})
    return outer


@api.get('/matrix', response_model=MatrixOut)
async def matrix(rows: int, columns: int, request: Request):
    """
    Returns a "rows x columns" matrix.
    """
    if not rows or not columns:
        raise HTTPException(status.HTTP_406_NOT_ACCEPTABLE,
                            detail="Must specify rows and columns query params.")

    with tracer.start_active_span("matrix-handler") as scope:
        # `start_active_span` leverages thread-local storage;
        # avoids the need to explicitly inject context into function signatures
        scope.span.set_tag("url", request.url)
        scope.span.log_kv({"src-ip": request.client.host, "rows": rows, "columns": columns})
        m = create_matrix(rows, columns)

    logger.info(f'returning a {rows}x{columns} matrix.')
    return MatrixOut(rows=rows, columns=columns, matrix=m)


@api.post('/classify')
async def classify():
    return {}
