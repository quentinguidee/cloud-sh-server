package com.quentinguidee.plugins

import com.quentinguidee.dao.oAuthMethodsDAO
import com.quentinguidee.models.*
import com.quentinguidee.utils.DatabaseConnectionFailedException
import io.ktor.server.application.*
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromStream
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction
import kotlin.io.path.Path

val DB_CONFIG_PATH = Path("data", "database.json")

fun Application.configureDatabase() {
    if (!DB_CONFIG_PATH.toFile().exists()) return

    connectDatabase()
    initDatabase()
}

val tables = arrayOf(
    Buckets,
    GitHubUsers,
    Nodes,
    Sessions,
    Users,
    UsersBuckets,
    UsersNodes,
    OAuthMethods,
)

@Serializable
data class DatabaseConfig(
    val host: String,
    val name: String,
    val user: String,
    var password: String,
)

@OptIn(ExperimentalSerializationApi::class)
fun getDatabaseConfig() = Json.decodeFromStream<DatabaseConfig>(
    DB_CONFIG_PATH.toFile().inputStream()
)

fun connectDatabase(config: DatabaseConfig = getDatabaseConfig()) {
    val database = Database.connect(
        "jdbc:postgresql://${config.host}/${config.name}",
        driver = "org.postgresql.Driver",
        user = config.user,
        password = config.password,
    )

    try {
        transaction(database) { connection.isClosed }
    } catch (e: Exception) {
        throw DatabaseConnectionFailedException()
    }
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
