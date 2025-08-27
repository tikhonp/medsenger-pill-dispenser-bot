<script setup>
import { computed, ref, watch } from 'vue'
import { usePillDispenserStore } from '@/stores/DispenserStore.js'

const props = defineProps({
  cols: {
    type: Number,
    required: true
  },
  rows: {
    type: Number,
    required: true
  }
})

const store = usePillDispenserStore()
const cells = computed(() => store.dispenserCells)

const step = ref(0)
const activeCell = ref(null)
const instruction = ref('Загрузка...')

function makeInstruction(i) {
  if (i >= cells.value.length) return 'Все таблетки разложены!'
  if (!cells.value[i].pills.length) return `Ячейка ${cells[i].cell_id} пустая`
  const pillsText = cells.value[i].pills.map(p => `<b>${p.name}</b> (${p.count} шт.)`).join(' | ')
  return `Положите в ячейку ${cells.value[i].cell_id}: \n${pillsText}`
}

function nextStep() {
  if (step.value < cells.value.length) {
    step.value++
    activeCell.value = cells.value[step.value]?.cell_id ?? null
    instruction.value = makeInstruction(step.value)
  }
}

function prevStep() {
  if (step.value > 0) {
    step.value--
    activeCell.value = cells.value[step.value]?.cell_id ?? null
    instruction.value = makeInstruction(step.value)
  }
}

watch(cells, (newCells) => {
  console.log(newCells)
  if (newCells.length > 0) {
    step.value = 0
    activeCell.value = newCells[0].cell_id
    instruction.value = makeInstruction(0)
  }
}, { immediate: true })
</script>

<template>

  <p class='instruction text-center' v-html='instruction'></p>

  <div class='d-flex justify-content-center gap-3'>
    <button class='btn btn-secondary arrow-btn left' @click='prevStep' :disabled='step === 0'>Назад</button>
    <button class='btn btn-secondary arrow-btn right' @click='nextStep' :disabled='step === cells.length'>Далее</button>
  </div>

  <div class='dispenser-wrapper'>
    <div class='dispenser-container' :style='{
        gridTemplateColumns: `repeat(${cols}, 150px)`,
        gridTemplateRows: `repeat(${rows}, 150px)`
      }'>
      <div v-for='cell in cells' :key='cell.cell_id'
           class='cell' :class='{ active: activeCell === cell.cell_id }'>
        <span class='badge rounded-pill text-primary cell-title'>{{ cell.cell_id }}</span>
        <ul class='pills-list'>
          <li v-for='pill in cell.pills' :key='pill.name'>
            <b>{{ pill.name }}</b> – {{ pill.count }} шт.
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dispenser-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 15px;
  border-radius: 12px;
  background-color: #b8e8ea;
  width: max-content;
  margin-left: auto;
  margin-right: auto;
}

.dispenser-container {
  display: grid;
  gap: 10px;
}

.cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border-radius: 8px;
  font-weight: bold;
  color: #f2f2f2;
  background-color: #24a8b4;
  transition: background-color 0.3s, transform 0.2s;
  padding: 6px;
  text-align: center;
  position: relative;
}

.cell.active {
  background-color: #006c88;
  border-color: #006c88;
}

.cell-title {
  position: absolute;
  top: 10px;
  margin-bottom: 5px;
  background-color: #b8e8ea;
  font-size: 13px;
  width: 2rem;
}


.pills-list {
  list-style: none;
  padding: 0;
  margin: 0;
  font-size: 0.85em;
}

.arrow-btn {
  margin: 0 5px;
}

.arrow-btn.left {
  border-radius: 0 8px 8px 0;
  clip-path: polygon(20px 0, 100% 0, 100% 100%, 20px 100%, 0 50%);
}

.arrow-btn.right {
  border-radius: 8px 0 0 8px;
  clip-path: polygon(0 0, calc(100% - 20px) 0, 100% 50%, calc(100% - 20px) 100%, 0 100%);
}

.instruction {
  white-space: pre-line;
}

.btn {
  width: 150px;
}
</style>