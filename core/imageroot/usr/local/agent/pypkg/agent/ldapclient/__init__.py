#
# Copyright (C) 2023 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License,
# or any later version.
#
# NethServer is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with NethServer.  If not, see COPYING.
#

from .exceptions import *
from .ad import LdapclientAd
from .rfc2307 import LdapclientRfc2307

class Ldapclient:
    @staticmethod
    def factory(**kwargs):
        if kwargs['schema'] == 'ad':
            return LdapclientAd(**kwargs)
        elif kwargs['schema'] == 'rfc2307':
            return LdapclientRfc2307(**kwargs)
        else:
            raise LdapclientUnknownSchema
