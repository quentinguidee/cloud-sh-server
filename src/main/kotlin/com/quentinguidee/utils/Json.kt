package com.quentinguidee.utils

import kotlinx.serialization.KSerializer
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonElement
import kotlinx.serialization.json.encodeToJsonElement
import kotlinx.serialization.json.jsonObject
import java.time.LocalDateTime

inline fun <reified T> json(
    base: T,
    addMore: MutableMap<String, JsonElement>.() -> Unit
): MutableMap<String, JsonElement> {
    val json = Json.encodeToJsonElement(base).jsonObject.toMutableMap()
    json.addMore()
    return json
}

inline fun <reified T> MutableMap<String, JsonElement>.putObject(key: String, value: T): JsonElement? {
    if (this.containsKey("${key}_id")) {
        this.remove("${key}_id")
    } else if (this.containsKey("${key}_uuid")) {
        this.remove("${key}_uuid")
    }
    return put(key, Json.encodeToJsonElement(value))
}

object DateSerializer : KSerializer<LocalDateTime> {
    override val descriptor = PrimitiveSerialDescriptor("LocalDateTime", PrimitiveKind.LONG)
    override fun serialize(encoder: Encoder, value: LocalDateTime) = encoder.encodeString(value.toString())
    override fun deserialize(decoder: Decoder): LocalDateTime = LocalDateTime.parse(decoder.decodeString())
}
