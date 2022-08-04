package com.quentinguidee.plugins

import com.quentinguidee.dao.oAuthMethodsDAO
import com.quentinguidee.models.*
import com.quentinguidee.utils.DatabaseConnectionFailedException
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromStream
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction
import kotlin.io.path.Path

val DB_CONFIG_PATH = Path("data", "database.json")

fun configureDatabase() {
    if (!DB_CONFIG_PATH.toFile().exists()) return

    connectDatabase()
}

val tables = arrayOf(
    Buckets,
    OAuthUsers,
    Nodes,
    Sessions,
    Users,
    UsersBuckets,
    UsersNodes,
    OAuthMethods,
)

@Serializable
data class DatabaseConfig(
    val dbms: String,
    val host: String? = null,
    val name: String? = null,
    val user: String? = null,
    var password: String? = null,
)

@OptIn(ExperimentalSerializationApi::class)
fun getDatabaseConfig() = Json.decodeFromStream<DatabaseConfig>(
    DB_CONFIG_PATH.toFile().inputStream()
)

fun connectDatabase(config: DatabaseConfig = getDatabaseConfig()) {
    val database = when (config.dbms) {
        "sqlite" -> Database.connect(
            "jdbc:${config.dbms}:data/database.db",
            driver = "org.sqlite.JDBC",
        )

        "postgresql" -> Database.connect(
            "jdbc:${config.dbms}://${config.host}/${config.name}",
            driver = "org.postgresql.Driver",
            user = config.user!!,
            password = config.password!!
        )

        else -> throw DatabaseConnectionFailedException()
    }

    try {
        transaction(database) { connection.isClosed }
    } catch (e: Exception) {
        throw DatabaseConnectionFailedException()
    }

    initDatabase()
}

fun initDatabase() = transaction {
    SchemaUtils.create(*tables)
}

fun resetDatabase() = transaction {
    val oAuthMethods = oAuthMethodsDAO.getAllPrivate()

    SchemaUtils.drop(*tables)
    initDatabase()

    oAuthMethods.forEach { oAuthMethodsDAO.create(it) }
}
