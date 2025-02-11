<template>
  <div>
    <div class="flex overflow-hidden border-solid border-0 border-b border-[#dfe0e5]">
      <!--CPU使用率-->
      <div class="flex-1 p-[20px] h-[360px] border-solid border-0 border-r border-[#dfe0e5]">
        <div class="flex justify-between">
          <span class="text-[14px] font-bold">{{ $t('CPU使用率') }}</span>
          <div>
            <div class="flex justify-end">
              <span class="text-[32px]">
                {{ conversionPercentUsed(overviewData.cpu_usage.used, overviewData.cpu_usage.total) }}
              </span>
              <sup class="text-[20px]">%</sup>
            </div>
            <div class="font-bold text-[#3ede78] text-[14px]">
              {{ parseFloat(overviewData.cpu_usage.used || 0).toFixed(2) }}
              of
              {{ parseFloat(overviewData.cpu_usage.total || 0).toFixed(2) }}
            </div>
          </div>
        </div>
        <ClusterOverviewChart
          :cluster-id="clusterId"
          class="!h-[250px]"
          :color="['#3ede78']"
          :metrics="['cpu_usage']"
        />
      </div>
      <!-- 内存使用率 -->
      <div class="flex-1 p-[20px] h-[360px] border-solid border-0 border-r border-[#dfe0e5]">
        <div class="flex justify-between">
          <span class="text-[14px] font-bold">{{ $t('内存使用率') }}</span>
          <div>
            <div class="flex justify-end">
              <span class="text-[32px]">
                {{ conversionPercentUsed(overviewData.memory_usage.used_bytes, overviewData.memory_usage.total_bytes) }}
              </span>
              <sup class="text-[20px]">%</sup>
            </div>
            <div class="font-bold text-[#3a84ff] text-[14px]">
              {{ formatBytes(overviewData.memory_usage.used_bytes || 0) }}
              of
              {{ formatBytes(overviewData.memory_usage.total_bytes || 0) }}
            </div>
          </div>
        </div>
        <ClusterOverviewChart
          :cluster-id="clusterId"
          class="!h-[250px]"
          :colors="['#3a84ff']"
          :metrics="['memory_usage']"
        />
      </div>
      <!-- 磁盘容量 -->
      <div class="flex-1 p-[20px] h-[360px]">
        <div class="flex justify-between">
          <span class="text-[14px] font-bold">{{ $t('磁盘容量') }}</span>
          <div>
            <div class="flex justify-end">
              <span class="text-[32px]">
                {{ conversionPercentUsed(overviewData.disk_usage.used_bytes, overviewData.disk_usage.total_bytes) }}
              </span>
              <sup class="text-[20px]">%</sup>
            </div>
            <div class="font-bold text-[#853cff] text-[14px]">
              {{ formatBytes(overviewData.disk_usage.used_bytes || 0) }}
              of
              {{ formatBytes(overviewData.disk_usage.total_bytes || 0) }}
            </div>
          </div>
        </div>
        <ClusterOverviewChart
          :cluster-id="clusterId"
          class="!h-[250px]"
          :colors="['#853cff']"
          :metrics="['disk_usage']"
        />
      </div>
    </div>
    <div class="flex overflow-hidden">
      <!--CPU装箱率-->
      <div class="flex-1 p-[20px] h-[360px] border-solid border-0 border-r border-[#dfe0e5]">
        <div class="flex justify-between">
          <span class="text-[14px] font-bold">{{ $t('CPU装箱率') }}</span>
          <div>
            <div class="flex justify-end">
              <span class="text-[32px]">
                {{ conversionPercentUsed(overviewData.cpu_usage.request, overviewData.cpu_usage.total) }}
              </span>
              <sup class="text-[20px]">%</sup>
            </div>
            <div class="font-bold text-[#3ede78] text-[14px]">
              {{ parseFloat(overviewData.cpu_usage.request || 0).toFixed(2) }}
              of
              {{ parseFloat(overviewData.cpu_usage.total || 0).toFixed(2) }}
            </div>
          </div>
        </div>
        <ClusterOverviewChart
          :cluster-id="clusterId"
          class="!h-[250px]"
          :color="['#3ede78']"
          :metrics="['cpu_request_usage']"
        />
      </div>
      <!-- 内存装箱率 -->
      <div class="flex-1 p-[20px] h-[360px] border-solid border-0 border-r border-[#dfe0e5]">
        <div class="flex justify-between">
          <span class="text-[14px] font-bold">{{ $t('内存装箱率') }}</span>
          <div>
            <div class="flex justify-end">
              <span class="text-[32px]">
                {{
                  conversionPercentUsed(
                    overviewData.memory_usage.request_bytes,
                    overviewData.memory_usage.total_bytes
                  )
                }}
              </span>
              <sup class="text-[20px]">%</sup>
            </div>
            <div class="font-bold text-[#3a84ff] text-[14px]">
              {{ formatBytes(overviewData.memory_usage.request_bytes || 0) }}
              of
              {{ formatBytes(overviewData.memory_usage.total_bytes || 0) }}
            </div>
          </div>
        </div>
        <ClusterOverviewChart
          :cluster-id="clusterId"
          class="!h-[250px]"
          :colors="['#3a84ff']"
          :metrics="['memory_request_usage']"
        />
      </div>
      <!-- 磁盘IO -->
      <div class="flex-1 p-[20px] h-[360px]">
        <div class="flex justify-between">
          <span class="text-[14px] font-bold">{{ $t('磁盘IO') }}</span>
          <div>
            <div class="flex justify-end">
              <span class="text-[32px]">
                {{ conversionPercentUsed(
                  overviewData.diskio_usage.used, overviewData.diskio_usage.total) }}
              </span>
              <sup class="text-[20px]">%</sup>
            </div>
            <!-- <div class="font-bold text-[#853cff] text-[14px]">
              {{ formatBytes(overviewData.disk_usage.used || 0) }}
              of
              {{ formatBytes(overviewData.disk_usage.total || 0) }}
            </div> -->
          </div>
        </div>
        <ClusterOverviewChart
          :cluster-id="clusterId"
          class="!h-[250px]"
          :colors="['#853cff']"
          :metrics="['diskio_usage']"
        />
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, onMounted, ref } from '@vue/composition-api';
import ClusterOverviewChart from './cluster-overview-chart.vue';
import $store from '@/store/index';
import { useProject } from '@/composables/use-app';
import { formatBytes } from '@/common/util';
export default defineComponent({
  name: 'ClusterOverview',
  components: { ClusterOverviewChart },
  setup(props, ctx) {
    const { $route } = ctx.root;
    const { projectCode } = useProject();
    const clusterId = computed(() => $route.params.clusterId);
    const overviewData = ref<{
      cpu_usage: any
      disk_usage: any
      memory_usage: any
      diskio_usage: any
    }>({
      cpu_usage: {},
      disk_usage: {},
      memory_usage: {},
      diskio_usage: {},
    });
    const getClusterOverview = async () => {
      overviewData.value = await $store.dispatch('metric/clusterOverview', {
        $projectCode: projectCode.value,
        $clusterId: clusterId.value,
      });
    };
    const conversionPercentUsed = (used, total) => {
      if (!total || parseFloat(total) === 0) {
        return 0;
      }
      let ret: any = parseFloat(used || 0) / parseFloat(total) * 100;
      if (ret !== 0 && ret !== 100) {
        ret = ret.toFixed(2);
      }
      return ret;
    };
    onMounted(() => {
      getClusterOverview();
    });
    return {
      overviewData,
      clusterId,
      conversionPercentUsed,
      formatBytes,
    };
  },
});
</script>
