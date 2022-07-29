package com.quentinguidee.services

import com.quentinguidee.dao.sessionsDAO
import org.jetbrains.exposed.sql.transactions.transaction

class SessionsServices {
    suspend fun session(token: String) = transaction {
        sessionsDAO.get(token)
    }

    suspend fun revokeSession(token: String) = transaction {
        sessionsDAO.delete(token)
    }

    suspend fun createSession(userID: Int) = transaction {
        sessionsDAO.create(userID)
    }
}

val sessionsServices = SessionsServices()
