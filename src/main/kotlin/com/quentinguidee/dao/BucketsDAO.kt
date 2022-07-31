package com.quentinguidee.dao

import com.quentinguidee.models.*
import org.jetbrains.exposed.sql.*
import java.util.*

class BucketsDAO {
    private fun toBucket(row: ResultRow) = Bucket(
        uuid = row[Buckets.id].value.toString(),
        name = row[Buckets.name],
        type = row[Buckets.type],
        size = row[Buckets.size],
        maxSize = row[Buckets.maxSize],
    )

    fun create(name: String, type: BucketType) = Buckets
        .insert {
            it[Buckets.name] = name
            it[Buckets.type] = type
        }.resultedValues!!.map(::toBucket).first()

    fun getUserBucket(userID: Int) = Buckets
        .innerJoin(UsersBuckets)
        .innerJoin(Users)
        .select { Users.id eq userID and (Buckets.type eq BucketType.USER_BUCKET) }
        .map(::toBucket)
        .first()

    fun increaseSize(bucketUUID: UUID, size: Int? = null) {
        if (size == null || size == 0) return
        Buckets.update({ Buckets.id eq bucketUUID }) {
            with(SqlExpressionBuilder) {
                it.update(Buckets.size, Buckets.size + size)
            }
        }
    }

    fun decreaseSize(bucketUUID: UUID, size: Int? = null) =
        increaseSize(bucketUUID, size?.let { -it })
}

val bucketsDAO = BucketsDAO()
