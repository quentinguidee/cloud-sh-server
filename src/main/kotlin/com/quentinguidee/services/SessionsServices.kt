package com.quentinguidee.services

import com.quentinguidee.dao.sessionsDAO
import com.quentinguidee.utils.NotAuthenticatedException
import org.jetbrains.exposed.sql.transactions.transaction

class SessionsServices {
    fun session(token: String) = transaction {
        try {
            sessionsDAO.get(token)
        } catch (e: NoSuchElementException) {
            throw NotAuthenticatedException()
        }
    }

    fun revokeSession(token: String) = transaction {
        sessionsDAO.delete(token)
    }

    fun createSession(userID: Int) = transaction {
        sessionsDAO.create(userID)
    }
}

val sessionsServices = SessionsServices()
