package com.quentinguidee.dao

import com.quentinguidee.models.db.*
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select

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
        .select {
            Users.id eq userID and
                    (Buckets.type eq BucketType.USER_BUCKET)
        }
        .map(::toBucket)
        .first()
}

val bucketsDAO = BucketsDAO()
