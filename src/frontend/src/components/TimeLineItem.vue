<script setup lang="ts" name="TimeLineItem">
import { showImagePreview } from 'vant'
import { computed } from 'vue'

// get props from parent

const props = defineProps({
  poster: String,
  timestamp: String,
  content: String,
  title: String,
  images: Array<string>,
})
let date = computed(() => {
  return props.timestamp?.slice(0, 10)
})

// auto set image columns
const images = props.images as Array<string>
let image_num = images ? images?.length : 0
let image_columns = image_num >= 3 ? 3 : image_num
console.log(`Image: ${image_num} col: ${image_columns}, img_num: ${image_num}`)

// open image preview
const openImagePreview = (index: number) => {
  console.log(`Image index: ${index}`)
  showImagePreview({
    images: props.images ? props.images : [], // write in this way to reslove type check warning.
    startPosition: index,
    closeable: true,
  })
}
</script>

<template>
  <div class="time-line-item">
    <!-- Time Block -->
    <div class="time-block">
      <p>{{ date }}</p>
    </div>

    <!-- Content Show -->
    <div class="content-block">
      <!-- Avatar & Name -->
      <div class="avatar-name">
        <van-image round width="1rem" height="1rem" :src="poster" />
        <span>{{ poster }}</span>
        <div class="title">{{ title }}</div>
      </div>

      <!-- Content -->
      <div class="content">
        <van-text-ellipsis
          :content="content"
          :rows="10"
          :expand-text="'展开'"
          :collapse-text="'收起'"
        ></van-text-ellipsis>

        <!-- Images -->
        <div>
          <van-grid
            square
            :v-show="image_num != 0"
            :column-num="image_columns"
            class="images-grid"
          >
            <van-grid-item v-for="(item, index) in images" :key="index">
              <van-image
                :src="item"
                fit="contain"
                width="100%"
                height="100%"
                position="center"
                @click="openImagePreview(index)"
                :preview-src="{ src: item }"
              />
            </van-grid-item>
          </van-grid>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.time-line-item {
  display: flex;
  background-color: gray;
  border-radius: 20px;
  margin: 5px 10px;

  .time-block {
    width: 100px;
    min-width: 70px;
    p {
      border: 1px solid #fff;
      border-radius: 20px;
    }
  }
  .content-block {
    display: flex;
    flex-direction: column;
    .avatar-name {
      display: flex;
      float: left;
      align-items: center;
      .title {
        color: black;
        margin-left: 10px;
        font-size: 8px;
        border-radius: 20px;
        background-color: aquamarine;
        padding: 4px;
      }
    }
    .content {
      .images-grid {
        margin: 5px;
      }
      .p {
        justify-content: left;
      }
      van-text-ellipsis {
        justify-content: left;
      }
    }
  }
}
</style>
