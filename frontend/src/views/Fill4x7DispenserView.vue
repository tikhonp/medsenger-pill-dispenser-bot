<script setup>
import FillDispenser from '@/components/FillDispenser.vue'
import Dispenser from '@/components/Dispenser.vue'
</script>

<template>
  <fill-dispenser layout='4x7'>
    <template #dispenser-type>размером 4 x 7</template>
    <template #dispenser>
      <dispenser :cols='7' :rows='4' />
    </template>
  </fill-dispenser>
</template>