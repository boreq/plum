<template>
  <div class="dashboard">
      <div class="parameters">
          <div class="box box-dimmed">
              <ul>
                  <li v-for="website of websites" :key="website">
                      <a v-bind:class="{ active: website === selectedWebsite }"
                         v-on:click="selectWebsite(website)">
                          {{ website }}
                      </a>
                  </li>
              </ul>
          </div>
      </div>

      <div class="parameters">
          <div class="box box-dimmed">
              <ul>
                  <li><a v-bind:class="{ active: selectedTimePeriod === TimePeriod.Day }" v-on:click="selectTimePeriod(TimePeriod.Day)">1D</a></li>
                  <li><a v-bind:class="{ active: selectedTimePeriod === TimePeriod.Week }" v-on:click="selectTimePeriod(TimePeriod.Week)">1W</a></li>
                  <li><a v-bind:class="{ active: selectedTimePeriod === TimePeriod.Month }" v-on:click="selectTimePeriod(TimePeriod.Month)">1M</a></li>
                  <li><a v-bind:class="{ active: selectedTimePeriod === TimePeriod.Year }" v-on:click="selectTimePeriod(TimePeriod.Year)">1Y</a></li>
              </ul>
          </div>

          <div class="box box-dimmed">
              <ul>
                  <li v-if="groupingTypeAvailable(GroupingType.Hourly)">
                      <a v-bind:class="{ active: selectedGroupingType === GroupingType.Hourly }"
                         v-on:click="selectGroupingType(GroupingType.Hourly)">
                          Hourly
                      </a>
                  </li>
                  <li v-if="groupingTypeAvailable(GroupingType.Daily)">
                      <a v-bind:class="{ active: selectedGroupingType === GroupingType.Daily }"
                         v-on:click="selectGroupingType(GroupingType.Daily)">
                          Daily
                      </a>
                  </li>
                  <li v-if="groupingTypeAvailable(GroupingType.Monthly)">
                      <a v-bind:class="{ active: selectedGroupingType === GroupingType.Monthly }"
                         v-on:click="selectGroupingType(GroupingType.Monthly)">
                          Monthly
                      </a>
                  </li>
              </ul>
          </div>

          <div class="box box-dimmed" v-for="activeFilter of activeFilters" :key="activeFilter.dimension">
              <a class="active"
                 title="Click to remove this filter."
                 v-on:click="removeFilter(activeFilter.dimension)">
                  {{ activeFilter.label }}: {{ activeFilter.value }}
              </a>
          </div>

          <div class="box box-dimmed" v-if="selectedRangeData">
              <a title="Click to cancel." v-on:click="selectData(null)">You are inspecting a single data point.</a>
          </div>

          <div class="box box-dimmed" v-if="updating || updatingLatest">
              <i class="fas fa-spinner fa-spin"></i>
          </div>
      </div>

      <div class="grid" :class="{ updating: updating }">
          <div class="box box-inversed malicious-traffic">
              <CategoryTraffic :data="data"
                               :category="Category.Malicious"
                               :checked="categoryChecked(Category.Malicious)"
                               :active="selectedCategory === Category.Malicious"
                               v-on:filter="toggleCategory"></CategoryTraffic>
          </div>
          <div class="box box-inversed automated-traffic">
              <CategoryTraffic :data="data"
                               :category="Category.Automated"
                               :checked="categoryChecked(Category.Automated)"
                               :active="selectedCategory === Category.Automated"
                               v-on:filter="toggleCategory"></CategoryTraffic>
          </div>
          <div class="box box-inversed possibly-automated-traffic">
              <CategoryTraffic :data="data"
                               :category="Category.PossiblyAutomated"
                               :checked="categoryChecked(Category.PossiblyAutomated)"
                               :active="selectedCategory === Category.PossiblyAutomated"
                               v-on:filter="toggleCategory"></CategoryTraffic>
          </div>
          <div class="box box-inversed unclassified-traffic">
              <CategoryTraffic :data="data"
                               :category="Category.Unclassified"
                               :checked="categoryChecked(Category.Unclassified)"
                               :active="selectedCategory === Category.Unclassified"
                               v-on:filter="toggleCategory"></CategoryTraffic>
          </div>
          <div class="box box-normal hits-and-visits">
              <HitsAndVisits :data="data" :groupingType="selectedGroupingType" v-on:select-data="selectData($event)"></HitsAndVisits>
          </div>
          <div class="box box-normal pages">
              <Pages :key="tableKey" :data="data" v-on:filter="addFilter"></Pages>
          </div>
          <div class="box box-normal referers">
              <Referers :key="tableKey" :data="data" v-on:filter="addFilter"></Referers>
          </div>
          <div class="box box-normal user-agents">
              <UserAgents :key="tableKey" :data="data" v-on:filter="addFilter"></UserAgents>
          </div>
          <div class="box box-normal bytes-sent-chart">
              <BytesSentChart :data="data" :groupingType="selectedGroupingType" v-on:select-data="selectData($event)"></BytesSentChart>
          </div>
          <div class="box box-normal bytes-sent">
              <BytesSent :key="tableKey" :data="data" v-on:filter="addFilter"></BytesSent>
          </div>
          <div class="box box-normal status-codes-chart">
              <StatusCodesChart :data="data" :groupingType="selectedGroupingType" v-on:select-data="selectData($event)"></StatusCodesChart>
          </div>
          <div class="box box-normal status-codes">
              <StatusCodes :key="tableKey" :data="data" v-on:filter="addFilter"></StatusCodes>
          </div>
      </div>
  </div>
</template>
<script lang="ts" src="./Dashboard.ts"></script>
<style scoped lang="scss" src="./Dashboard.scss"></style>
