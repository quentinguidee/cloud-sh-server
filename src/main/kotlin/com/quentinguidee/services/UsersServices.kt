package com.quentinguidee.services

import com.quentinguidee.models.db.User
import com.quentinguidee.models.db.Users
import org.jetbrains.exposed.sql.transactions.transaction

class UsersServices {
    suspend fun get(username: String) = transaction {
        User
            .find { Users.username eq username }
            .single()
    }
}

val usersServices = UsersServices()
