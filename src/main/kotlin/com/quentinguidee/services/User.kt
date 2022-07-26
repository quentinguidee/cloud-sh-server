package com.quentinguidee.services

import com.quentinguidee.models.User
import com.quentinguidee.models.Users
import org.jetbrains.exposed.sql.transactions.transaction

class UserService {
    fun getUser(username: String): User? {
        return transaction {
            return@transaction User
                .find { Users.username eq username }
                .firstOrNull()
        }
    }

}

val userService = UserService()
