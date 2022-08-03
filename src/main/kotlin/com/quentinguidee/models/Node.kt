package com.quentinguidee.models

import com.quentinguidee.utils.DateSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.UUIDTable
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime
import kotlin.io.path.Path
import kotlin.io.path.extension

object Nodes : UUIDTable("nodes") {
    val parent = reference("parent_uuid", Nodes).nullable()
    val bucket = reference("bucket_uuid", Buckets)
    val name = varchar("name", 255)
    val type = varchar("type", 255)
    val mime = varchar("mime", 255).nullable()
    val size = integer("size").default(0)

    val createdAt = datetime("created_at").default(LocalDateTime.now())
    val updatedAt = datetime("updated_at").nullable()
    val deletedAt = datetime("deleted_at").nullable()
}

@Serializable
data class Node(
    val uuid: String,
    @SerialName("parent_uuid")
    val parentUUID: String?,
    @SerialName("bucket_uuid")
    val bucketUUID: String,
    val name: String,
    val type: String,
    val mime: String?,
    val size: Int,

    @Serializable(DateSerializer::class) val createdAt: LocalDateTime,
    @Serializable(DateSerializer::class) val updatedAt: LocalDateTime? = null,
    @Serializable(DateSerializer::class) val deletedAt: LocalDateTime? = null,
)

fun deduceNodeTypeByName(name: String) = when (Path(name).extension) {
    "afdesign" -> "afdesign"
    "afphoto" -> "afphoto"
    "afpub" -> "afpub"
    "avi" -> "avi"
    "babelrc" -> "babel"
    "bmp" -> "bmp"
    "c" -> "c"
    "cpp", "cxx" -> "cpp"
    "css" -> "css"
    "word", "odt", "doc", "docx" -> "document"
    "flac" -> "flac"
    "gitignore", "gitkeep" -> "git"
    "go" -> "go"
    "html", "htm" -> "html"
    "js" -> "javascript"
    "jpeg" -> "jpeg"
    "jpg" -> "jpg"
    "json" -> "json"
    "kt" -> "kotlin"
    "md" -> "markdown"
    "mkv" -> "mkv"
    "mov" -> "mov"
    "mp3" -> "mp3"
    "mp4" -> "mp4"
    "ml", "mli" -> "ocaml"
    "ogg" -> "ogg"
    "pdf" -> "pdf"
    "php" -> "php"
    "png" -> "png"
    "ppt", "pptx", "odp" -> "presentation"
    "py" -> "python"
    "raw" -> "raw"
    "tsx" -> "react"
    "rb" -> "ruby"
    "sass", "scss" -> "sass"
    "sc" -> "scala"
    "sh" -> "shell"
    "xls", "xlsx", "ods" -> "spreadsheet"
    "ts" -> "typescript"
    "wav" -> "wav"
    "lock" -> "yarn"
    else -> "file"
}
