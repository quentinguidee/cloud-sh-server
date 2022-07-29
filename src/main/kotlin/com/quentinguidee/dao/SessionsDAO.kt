package com.quentinguidee.dao

import com.quentinguidee.models.db.Session
import com.quentinguidee.models.db.Sessions
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import java.util.*

class SessionsDAO {
    private fun toSession(row: ResultRow) = Session(
        id = row[Sessions.id].value,
        userID = row[Sessions.user].value,
        token = row[Sessions.token],
    )

    fun get(userID: Int) = Sessions
        .select { Sessions.user eq userID }
        .map(::toSession)
        .first()


    fun get(token: String) = Sessions
        .select { Sessions.token eq token }
        .map(::toSession)
        .first()


    fun create(userID: Int) = Sessions.insert {
        it[user] = userID
        it[token] = UUID.randomUUID().toString()
    }.resultedValues?.map(::toSession)!!.first()

    fun delete(token: String) = Sessions
        .deleteWhere { Sessions.token eq token }
}

val sessionsDAO = SessionsDAO()
