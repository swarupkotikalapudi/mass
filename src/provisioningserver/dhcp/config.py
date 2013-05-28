# Copyright 2012 Canonical Ltd.  This software is licensed under the
# GNU Affero General Public License version 3 (see the file LICENSE).

"""Write config output for ISC DHCPD."""

from __future__ import (
    absolute_import,
    print_function,
    unicode_literals,
    )

__metaclass__ = type
__all__ = [
    "DHCPConfigError",
    "get_config",
]


from os import path
from platform import linux_distribution

from provisioningserver.pxe.tftppath import compose_bootloader_path
import tempita

# TODO: make this configurable.
template_dir = path.join(path.dirname(__file__), 'templates')


class DHCPConfigError(Exception):
    """Exception raised for errors processing the DHCP config."""


def get_config(**params):
    """Return a DHCP config file based on the supplied parameters.

    :param subnet: The base subnet declaration. e.g. 192.168.1.0
    :param subnet_mask: The mask for the above subnet, e.g. 255.255.255.0
    :param broadcast_address: The broadcast IP address for the subnet,
        e.g. 192.168.1.255
    :param dns_servers: One or more IP addresses of the DNS server for the
        subnet
    :param domain_name: name that will be appended to the client's hostname to
        form a fully-qualified domain-name
    :param gateway: The router/gateway IP address for the subnet.
    :param low_range: The first IP address in the range of IP addresses to
        allocate
    :param high_range: The last IP address in the range of IP addresses to
        allocate
    """
    template_file = path.join(template_dir, 'dhcpd.conf.template')
    params['bootloader'] = compose_bootloader_path()
    params['platform_codename'] = linux_distribution()[2]
    try:
        template = tempita.Template.from_filename(
            template_file, encoding="UTF-8")
        return template.substitute(params)
    except NameError, error:
        raise DHCPConfigError(*error.args)
