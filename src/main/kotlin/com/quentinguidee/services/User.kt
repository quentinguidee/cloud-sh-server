package com.quentinguidee.services

import com.quentinguidee.models.User
import com.quentinguidee.models.Users
import org.jetbrains.exposed.sql.transactions.transaction

class UserService {
    suspend fun get(username: String) = transaction {
        User
            .find { Users.username eq username }
            .singleOrNull()
    }
}

val userService = UserService()
