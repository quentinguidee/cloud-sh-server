package com.quentinguidee.utils

import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonElement
import kotlinx.serialization.json.encodeToJsonElement
import kotlinx.serialization.json.jsonObject

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
