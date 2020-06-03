/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package io.charlescd.moove.commons.representation

import com.fasterxml.jackson.annotation.JsonFormat
import java.time.LocalDateTime

data class CardRepresentation(
    val id: String,
    val name: String,
    val description: String?,
    val column: CardColumnRepresentation,
    val author: SimpleUserRepresentation,
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd HH:mm:ss")
    val createdAt: LocalDateTime,
    val labels: List<SimpleLabelRepresentation>,
    val type: String,
    val feature: FeatureRepresentation?,
    val hypothesisId: String,
    val comments: List<CommentRepresentation> = emptyList(),
    val members: List<UserRepresentation> = emptyList(),
    val index: Int?
)
