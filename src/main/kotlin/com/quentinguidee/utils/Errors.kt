package com.quentinguidee.utils

import com.quentinguidee.models.User

class DatabaseConnectionFailedException() :
    Exception("Failed to connect to the database.")

class ServerAlreadyConfiguredException() :
    Exception("The server is already configured. Use the admin routes to edit them, or delete your data/database.json config file.")

class NotAuthenticatedException :
    Exception("You're not logged in.")

class UnauthorizedException(user: User) :
    Exception("${user.username} is not authorized to access this")

class MaxRecursionLevelException :
    Exception("max recursion level reached")
