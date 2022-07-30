package com.quentinguidee.utils

import com.quentinguidee.models.User

class UnauthorizedException(user: User) : Exception("${user.username} is not authorized to access this")
class MaxRecursionLevelException : Exception("max recursion level reached")
