<script setup lang="ts" name="TimeLineItem">
import { showImagePreview } from 'vant'

// get props from parent

const props = defineProps({
  poster: String,
  date: String,
  title: String,
  images: Array<string>,
})

// auto set image columns
const images = props.images as Array<string>
let image_num = images ? images?.length : 0
let image_columns = image_num >= 3 ? 3 : image_num
console.log(`Image: ${image_num} col: ${image_columns}`)

// open image preview
const openImagePreview = () => {
  showImagePreview({
    images: props.images ? props.images : [], // write in this way to reslove type check warning.
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
        <p>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas,
          doloremque.
        </p>
        <!-- Images -->
        <van-grid
          square
          :v-show="image_num != 0"
          :column-num="image_columns"
          class="images-grid"
        >
          <van-grid-item v-for="(item, index) in images" :key="index">
            <van-image
              :src="item"
              fit="cover"
              width="6rem"
              height="6rem"
              position="center"
              @click="openImagePreview"
              :preview-src="{ src: item }"
            />
          </van-grid-item>
        </van-grid>
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
    }
  }
}
</style>
