import { defineStore } from 'pinia'
import api from '@/api/agent_api.js'

export const usePillDispenserStore = defineStore('pillDispenser', {
  state: () => ({
    layout: null,
    loading: false,
    error: null
  }),

  getters: {
    pills(state) {
      const pills = {}
      state.layout?.cells?.forEach(cell => {
        pills[cell.cell_name] = (pills[cell.cell_name] || 0) + 1
      })
      console.log(pills)
      return pills
    },
    dispenserCells: (state) => {
      if (!state.layout) return []

      const hardwareConfig = {
        HW_2X2_V1: { rows: 2, cols: 2 },
        HW_4X7_V1: { rows: 4, cols: 7 }
      }

      const config = hardwareConfig[state.layout.hardware_type]
      if (!config) return []

      const totalCells = config.rows * config.cols

      const grouped = {}
      state.layout.cells.forEach(c => {
        if (!grouped[c.cell_id]) grouped[c.cell_id] = []
        grouped[c.cell_id].push({
          name: c.cell_name,
          time: c.time
        })
      })

      return Array.from({ length: totalCells }, (_, i) => ({
        cell_id: i + 1,
        pills: grouped[i + 1]?.map(p => ({ name: p.name, count: 1 })) || []
      }))
    }
  },

  actions: {
    async loadLayout(dispenser_type) {
      this.loading = true
      this.layout = null
      let layout = dispenser_type === '2x2' ? await api.get2x2Layout() : await api.getLargeLayout()

      if (layout !== null) {
        this.layout = layout
      } else {
        this.error = 'Ошибка загрузки.'
      }
      console.log(this.layout)

      this.loading = false
    }
  }
})