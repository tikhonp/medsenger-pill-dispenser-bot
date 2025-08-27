<script setup>
import { useRouter } from 'vue-router'
import { usePillDispenserStore } from '../stores/DispenserStore.js'
import { onMounted } from 'vue'

const router = useRouter()
const store = usePillDispenserStore()

function goBack() {
  router.back()
}

const props = defineProps({
  layout: {
    type: String,
    required: true
  }
})

onMounted(() => {
  store.loadLayout(props.layout)
})

</script>

<template>
  <div>
    <div class='d-flex justify-content-end mb-3'>
      <button class='btn btn-light btn-sm float-end' @click='goBack'>
        Назад
      </button>
    </div>

    <!-- инструкция -->
    <h4 class='mb-3'>Подготовка</h4>
    <ol class='text-left mb-4'>
      <li class='mb-2'>
        Освободите <strong>таблетницу
        <slot name='dispenser-type'></slot>
        </strong> и расположите ее кнопкой на себя.
      </li>
      <li class='mb-2'>
        Подготовьте необходимые препараты:
        <div v-if='store.loading'>Загрузка...</div>
        <div class='alert alert-warning' v-else-if='store.error'>{{ store.error }}</div>
        <ul class='list-group mt-2' v-else>
          <li v-for='(cnt, pill) in store.pills' :key='pill'
              class='list-group-item d-flex justify-content-between align-items-center'>
            <span>{{ pill }}</span>
            <span class='badge pill-counter rounded-pill text-primary'>{{ cnt }} шт</span>
          </li>
        </ul>
      </li>
      <li class='mb-4'>
        Заполните ячейки, следуя инструкциям ниже.
      </li>
    </ol>
    <slot name='dispenser'></slot>

  </div>
</template>

<style scoped>
main {
  display: flex;
  flex-direction: column;
}

.pill-counter {
  background-color: #b8e8ea;
  font-size: 13px;
  width: 4rem;
}
</style>