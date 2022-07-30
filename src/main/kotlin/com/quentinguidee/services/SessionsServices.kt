package com.quentinguidee.services

import com.quentinguidee.dao.sessionsDAO
import org.jetbrains.exposed.sql.transactions.transaction

class SessionsServices {
    fun session(token: String) = transaction {
        sessionsDAO.get(token)
    }

    fun revokeSession(token: String) = transaction {
        sessionsDAO.delete(token)
    }

    fun createSession(userID: Int) = transaction {
        sessionsDAO.create(userID)
    }
}

val sessionsServices = SessionsServices()
