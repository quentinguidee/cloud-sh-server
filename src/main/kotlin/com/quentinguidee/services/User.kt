package com.quentinguidee.services

import com.quentinguidee.dao.daoUser
import org.jetbrains.exposed.sql.transactions.transaction

class UserService {
    suspend fun get(username: String) = transaction {
        daoUser.get(username)
    }
}

val userService = UserService()
