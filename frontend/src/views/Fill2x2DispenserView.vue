<script setup>
import FillDispenser from '@/components/FillDispenser.vue'
import Dispenser from '@/components/Dispenser.vue'
</script>

<template>
  <fill-dispenser layout='2x2'>
    <template #dispenser-type>размером 2 x 2</template>
    <template #dispenser>
      <dispenser :cols='2' :rows='2' />
    </template>
  </fill-dispenser>
</template>
