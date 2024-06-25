from fastapi import Depends, Response

from maasapiserver.common.api.base import Handler, handler
from maasapiserver.common.api.models.responses.errors import (
    ValidationErrorBodyResponse,
)
from maasapiserver.v3.api import services
from maasapiserver.v3.api.models.requests.query import TokenPaginationParams
from maasapiserver.v3.api.models.responses.machines import MachinesListResponse
from maasapiserver.v3.auth.base import check_permissions
from maasapiserver.v3.auth.jwt import UserRole
from maasapiserver.v3.constants import V3_API_PREFIX
from maasapiserver.v3.services import ServiceCollectionV3


class MachinesHandler(Handler):
    """Machines API handler."""

    TAGS = ["Machines"]

    @handler(
        path="/machines",
        methods=["GET"],
        tags=TAGS,
        responses={
            200: {
                "model": MachinesListResponse,
            },
            422: {"model": ValidationErrorBodyResponse},
        },
        response_model_exclude_none=True,
        status_code=200,
        dependencies=[
            Depends(check_permissions(required_roles={UserRole.USER}))
        ],
    )
    async def list_machines(
        self,
        token_pagination_params: TokenPaginationParams = Depends(),
        services: ServiceCollectionV3 = Depends(services),
    ) -> Response:
        machines = await services.machines.list(
            token=token_pagination_params.token,
            size=token_pagination_params.size,
        )
        return MachinesListResponse(
            items=[
                machine.to_response(f"{V3_API_PREFIX}/machines")
                for machine in machines.items
            ],
            next=f"{V3_API_PREFIX}/machines?"
            f"{TokenPaginationParams.to_href_format(machines.next_token, token_pagination_params.size)}"
            if machines.next_token
            else None,
        )
