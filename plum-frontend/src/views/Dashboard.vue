<template>
  <div class="dashboard">
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

          <div class="box box-dimmed">
              <ul>
                  <li v-for="website of websites">
                      <a v-bind:class="{ active: website === selectedWebsite }"
                         v-on:click="selectWebsite(website)">
                          {{ website }}
                      </a>
                  </li>
              </ul>
          </div>

          <div class="box box-dimmed" v-if="selectedRangeData">
              <a title="Click to cancel." v-on:click="selectData(null)">You are inspecting a single data point.</a>
          </div>

          <div class="box box-dimmed" v-if="updating || updatingLatest">
              <i class="fas fa-spinner fa-spin"></i>
          </div>
      </div>

      <div class="grid" :class="{ updating: updating }">
          <div class="box box-inversed summary">
              <Summary :data="data"></Summary>
          </div>
          <div class="box box-normal hits-and-visits">
              <HitsAndVisits :data="data" :groupingType="selectedGroupingType" v-on:select-data="selectData($event)"></HitsAndVisits>
          </div>
          <div class="box box-normal pages">
              <Pages :data="data"></Pages>
          </div>
          <div class="box box-normal referers">
              <Referers :data="data"></Referers>
          </div>
          <div class="box box-normal bytes-sent-chart">
              <BytesSentChart :data="data" :groupingType="selectedGroupingType" v-on:select-data="selectData($event)"></BytesSentChart>
          </div>
          <div class="box box-normal bytes-sent">
              <BytesSent :data="data"></BytesSent>
          </div>
          <div class="box box-normal status-codes-chart">
              <StatusCodesChart :data="data" :groupingType="selectedGroupingType" v-on:select-data="selectData($event)"></StatusCodesChart>
          </div>
          <div class="box box-normal status-codes">
              <StatusCodes :data="data"></StatusCodes>
          </div>
      </div>
  </div>
</template>
<script lang="ts" src="./Dashboard.ts"></script>
<style scoped lang="scss" src="./Dashboard.scss"></style>
