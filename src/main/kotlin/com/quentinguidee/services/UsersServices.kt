package com.quentinguidee.services

import com.quentinguidee.dao.usersDAO
import org.jetbrains.exposed.sql.transactions.transaction

class UsersServices {
    suspend fun get(username: String) = transaction {
        usersDAO.get(username)
    }
}

val usersServices = UsersServices()
