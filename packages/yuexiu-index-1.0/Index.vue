<template>
	<v-row no-gutters="" align="start" id="homeMain">
		<!-- <div style="z-index:99999;position:fixed;" v-if="showHelp"
			><v-btn
				@click="showHelp = false"
				fixed=""
				fab
				top
				right
				color="lime"
				x-small=""
				elevation="1"
				class="black--text"
			>
				<v-icon dense="">fa-times</v-icon>
			</v-btn>
			<img
				src="../assets/help.png"
				:style="
					`height:${$store.state.windowHeight}px;width:${$store.state.windowWidth}px;`
				"
			/>
		</div> -->

		<v-col cols="12" class="bg">
			<v-app-bar color="rgba(0, 0, 0, 0.4)" dark="" flat="" dense="">
				<v-toolbar-title class="subtitle-1"
					><v-icon left="" dense="" color="cyan">fa-map-marker-alt</v-icon
					>{{ $t($store.state.startPointInfo.name) }}</v-toolbar-title
				>
				<v-spacer></v-spacer>

				<v-menu offset-y bottom>
					<template v-slot:activator="{ on }">
						<v-btn v-on="on" icon>
							<v-icon>fa-language</v-icon>
						</v-btn>
					</template>
					<v-list dense="">
						<v-list-item
							v-for="(item, index) in $i18n.messages"
							:key="index"
							@click="changeLang(index)"
							class=" subtitle-2"
						>
							{{ $t(index) }}
						</v-list-item>
					</v-list>
				</v-menu>
			</v-app-bar>
			<v-row
				no-gutters=""
				justify="center"
				align="end"
				style="margin-top: 12.7em;"
			>
				<div
					class="roundBorder"
					:style="`width:100%;background: rgba(0, 0, 0, 0.4)`"
				>
					<div
						class="roundBorder1"
						:style="`width:100%;border-top: 2px solid ${separatorColor};`"
					>
						<v-bottom-navigation
							background-color="rgba(0, 0, 0, 0)"
							grow
							color="teal"
							class=" elevation-0"
							dark
							height="80"
						>
							<template v-for="(item, index) in bottomNav1">
								<v-divider
									:key="`divider-${index}`"
									v-if="index > 0"
									inset
									vertical
									class="mt-4"
									:style="`border-color: ${separatorColor} !important`"
								></v-divider>
								<v-btn :key="`btn-${item.text}`" @click="navPage(item.text)">
									<span :class="`white--text mb-1 ${navFont}`">{{
										$t(item.text)
									}}</span>
									<v-img :src="item.image" width="60" contain=""></v-img>
								</v-btn>
							</template>
						</v-bottom-navigation> </div
				></div>
				<div :style="`width:100%;background: rgba(0, 0, 0, 0.5)`">
					<v-divider
						dark=""
						:style="`border-color: ${separatorColor} !important`"
					></v-divider>
					<v-bottom-navigation
						background-color="rgba(0, 0, 0, 0)"
						dark=""
						grow
						color="teal"
						class=" elevation-0 mb-1"
						height="80"
					>
						<template v-for="(item, index) in bottomNav2">
							<v-divider
								:key="`divider-${index}`"
								v-if="index > 0"
								inset
								vertical
								class="mt-0"
								:style="`border-color: ${separatorColor} !important`"
							></v-divider>
							<v-btn :key="`btn-${item.text}`" @click="navPage(item.text)">
								<span :class="`white--text ${navFont}`">{{
									$t(item.text)
								}}</span>
								<v-img :src="item.image" width="60" contain=""></v-img>
							</v-btn>
						</template> </v-bottom-navigation
				></div>
			</v-row>
		</v-col>
		<v-col cols="12" class="mt-0">
			<v-card elevation="0" tile="">
				<v-container>
					<v-card elevation="0" class="mb-3" tile="">
						<v-card-title class=" pl-0"
							><img
								width="130"
								:src="require(`@/assets/title-attractions_${$i18n.locale}.png`)"
						/></v-card-title>
						<swiper :options="swiperOption">
							<swiper-slide v-for="(b, i) in banner" :key="i"
								><v-img
									:src="b.img"
									aspect-ratio="1.5"
									@click="bannerClick(b)"
								></v-img
							></swiper-slide>
						</swiper>
					</v-card>
					<v-divider style="border-top:1px dashed #ccc"></v-divider>
					<v-card elevation="0">
						<v-card-title class=" pl-0"
							><img
								:src="require(`@/assets/title-news_${$i18n.locale}.png`)"
								width="130"
						/></v-card-title>
						<v-list three-line>
							<v-list-item
								v-for="item in articles"
								:key="item.article_id"
								@click="navPage('article/' + item.article_id)"
								class="mt-0"
							>
								<v-list-item-content class="pt-0">
									<v-list-item-subtitle
										class=" subtitle-2 black--text"
										v-text="item.article_title"
									></v-list-item-subtitle>
									<v-list-item-title class="overline grey--text">{{
										item.date
									}}</v-list-item-title>
								</v-list-item-content>

								<v-list-item-action v-if="item.article_icon">
									<v-img
										:src="`https://o.signp.cn/${item.article_icon}`"
									></v-img>
								</v-list-item-action>
							</v-list-item>
						</v-list>
					</v-card>
					<v-divider style="border-top:1px dashed #ccc"></v-divider>
					<v-card elevation="0">
						<v-card-title class=" pl-0"
							><img
								:src="require(`@/assets/title-weather_${$i18n.locale}.png`)"
								width="130"
						/></v-card-title>
						<v-row no-gutters="">
							<v-col cols="3" class="text-center">
								<v-card
									elevation="0"
									height="100%"
									class=" caption blue lighten-4"
								>
									{{ getWeek() }}<br /><span class="overline">{{
										dateFormat()
									}}</span>
								</v-card>
							</v-col>
							<v-col cols="9" class="pl-1 white">
								<v-row no-gutters="" align="start" justify="start"
									><v-col cols="12"
										><span class="subtitle-1 text--darken-1 grey--text">
											{{ weatherString }}
										</span></v-col
									></v-row
								>
							</v-col></v-row
						>
					</v-card>
				</v-container>
			</v-card>
		</v-col>
	</v-row>
</template>

<script>
import "swiper/dist/css/swiper.css";

import { swiper, swiperSlide } from "vue-awesome-swiper";
export default {
	components: { swiper, swiperSlide },
	mounted() {
		const _this = this;

		if (this.$i18n.locale === "zh_CN") {
			this.navFont = "subtitle-2";
		}
		_this.loadBanner();
		// const showHelp = localStorage.getItem("signp-145:showHelp");
		// if (!showHelp) this.showHelp = true;
	},
	computed: {
		locale() {
			return this.$i18n.locale;
		},
		bannerHeight() {
			return parseInt(this.$store.state.windowHeight / 4);
		},
		windowHeight() {
			return this.$store.state.windowHeight;
		}
	},
	watch: {
		locale: {
			handler(val) {
				let lang = val;
				if (!lang) lang = "zh_CN";
				if (lang === "zh_CN") {
					this.navFont = "subtitle-2";
				} else {
					this.navFont = "";
				}
				this.loadArticles();
				this.loadWeather();
				// this.loadPolygons();
			},
			immediate: true
		},
		showHelp(val) {
			if (val === true) {
				localStorage.setItem("signp-145:showHelp", "ok");
			}
		}
	},
	data() {
		return {
			swiperOption: {
				slidesPerView: 3,
				spaceBetween: 10,
				freeMode: true,
				pagination: {
					el: ".swiper-pagination",
					clickable: true
				}
			},
			showArticleNumber: 3,
			separatorColor: "rgba(255,255,255,0.3)",
			showHelp: false,
			bannerType: "round",
			articles: [],
			navFont: "",
			enterMap: false,
			transition: "scale-transition",
			markerOffset: [0, 0],
			markerList: [],
			weatherString: "",
			bottomNav1: [
				{
					text: "attractions",
					image: require("@/assets/v2-b1.png")
				},
				{
					text: "recommend",
					image: require("@/assets/v2-b6.png")
				},
				{
					text: "map",
					image: require("@/assets/v2-b3.png")
				},
				{
					text: "help",
					image: require("@/assets/v2-b8.png")
				}
			],
			bottomNav2: [
				{
					text: "about",
					image: require("@/assets/v2-b4.png")
				},
				{
					text: "surrounding",
					image: require("@/assets/v2-b7.png")
				},
				{
					text: "vr",
					image: require("@/assets/v2-b2.png")
				},
				{
					text: "search",
					image: require("@/assets/v2-b5.png")
				}
			],
			banner: []
		};
	},
	methods: {
		mapLimitBounds() {
			const bounds = new window.AMap.Bounds(
				new window.AMap.LngLat(
					this.$store.state.mapBounds.southWest[0],
					this.$store.state.mapBounds.southWest[1]
				),
				new window.AMap.LngLat(
					this.$store.state.mapBounds.northEast[0],
					this.$store.state.mapBounds.northEast[1]
				)
			);
			this.$refs.map.$$getInstance().setLimitBounds(bounds);
		},
		getWeek() {
			const weekDay = [
				this.$t("week7"),
				this.$t("week1"),
				this.$t("week2"),
				this.$t("week3"),
				this.$t("week4"),
				this.$t("week5"),
				this.$t("week6")
			];
			const d = new Date();
			return weekDay[d.getDay()];
		},
		dateFormat(timestamp, sep = "-") {
			let d = new Date();
			if (timestamp) {
				d = new Date(timestamp);
			}
			const nowDay = d.getDate();
			const day = nowDay >= 10 ? nowDay.toString() : "0" + nowDay.toString();
			const html = [
				d.getFullYear().toString(),
				(d.getMonth() + 1).toString(),
				day
			];
			return html.join(sep);
		},
		loadArticles() {
			const _this = this;
			_this.articles = [];
			const readStore = _this.$store.state.db
				.transaction("articles")
				.objectStore("articles")
				.get(
					_this.$store.state.startPointInfo.index_article_category +
						"_" +
						_this.$i18n.locale
				);
			readStore.onsuccess = function(e) {
				const r = e.target.result;
				const nowTime = new Date().getTime();
				if (r && r.timestamp > nowTime && r.val.length > 0) {
					for (let i = 0; i < r.val.length; i++) {
						if (i === _this.$store.state.startPointInfo.showArticleLimit) break;
						let _d = new Date(parseInt(r.val[i].article_create_at) * 1000);
						r.val[i].date =
							_d.getFullYear() + "/" + (_d.getMonth() + 1) + "/" + _d.getDate();
						_this.articles.push(r.val[i]);
					}
				} else {
					_this.$http
						.get(_this.apiHost + "/articles", {
							params: {
								categoryID:
									_this.$store.state.startPointInfo.index_article_category,
								locale: _this.$i18n.locale,
								limit: _this.$store.state.startPointInfo.articleLimit,
								projectID: _this.$store.state.startPointInfo.project_id,
								timestamp: parseInt(new Date().getTime() / 1000)
							}
						})
						.then(function(resp) {
							if (resp.status === 200) {
								const timestamp = new Date().getTime() + 300 * 1000;
								const writeStore = _this.$store.state.db
									.transaction("articles", "readwrite")
									.objectStore("articles");
								writeStore.put({
									article_category_id:
										_this.$store.state.startPointInfo.index_article_category +
										"_" +
										_this.$i18n.locale,
									timestamp: timestamp,
									val: resp.body
								});
								for (let i = 0; i < resp.body.length; i++) {
									if (i === _this.$store.state.startPointInfo.showArticleLimit)
										break;
									let _d = new Date(
										parseInt(resp.body[i].article_create_at) * 1000
									);
									resp.body[i].date =
										_d.getFullYear() +
										"/" +
										(_d.getMonth() + 1) +
										"/" +
										_d.getDate();
									_this.articles.push(resp.body[i]);
								}
							}
						});
				}
			};
		},
		loadBanner() {
			const _this = this;
			_this.articles = [];
			const readStore = _this.$store.state.db
				.transaction("banner")
				.objectStore("banner")
				.get(_this.bannerType);
			readStore.onsuccess = function(e) {
				const banner = [
					// {
					// 	img: require("@/assets/banner1.jpg"),
					// 	title: "",
					// 	link: "about",
					// 	mapGid: ""
					// }
				];
				const r = e.target.result;
				const nowTime = new Date().getTime();
				if (r && r.timestamp > nowTime && r.val.length > 0) {
					for (let i = 0; i < r.val.length; i++) {
						banner.push({
							img: "https://o.signp.cn/" + r.val[i].images_path,
							title: "",
							bindType: r.val[i].images_bind_type,
							link: r.val[i].images_link,
							mapGid: r.val[i].images_map_gid
						});
					}
					_this.banner = banner;
				} else {
					_this.$http
						.get(_this.apiHost + "/images", {
							params: {
								type: _this.bannerType,
								projectID: _this.$store.state.startPointInfo.project_id,
								timestamp: parseInt(nowTime / 1000)
							}
						})
						.then(function(resp) {
							if (resp.status === 200) {
								const writeStore = _this.$store.state.db
									.transaction("banner", "readwrite")
									.objectStore("banner");
								writeStore.put({
									image_type: _this.bannerType,
									timestamp: nowTime + 300 * 1000,
									val: resp.body
								});
								for (let i = 0; i < resp.body.length; i++) {
									banner.push({
										img: "https://o.signp.cn/" + resp.body[i].images_path,
										title: "",
										bindType: resp.body[i].images_bind_type,
										link: resp.body[i].images_link,
										mapGid: resp.body[i].images_map_gid
									});
								}
								_this.banner = banner;
							}
						});
				}
			};
		},
		loadWeather() {
			const weatherTag = "signp-145:weather-" + this.locale;
			const weather = JSON.parse(localStorage.getItem(weatherTag));
			const nowTime = new Date().getTime();
			if (!weather || weather.timestamp < nowTime) {
				this.$http
					.get(this.apiHost + "/weather", {
						params: {
							location: this.$store.state.startPointInfo.weatherLocation,
							lang: this.locale.substring(0, 2)
						}
					})
					.then(
						function(resp) {
							if (resp.status === 200) {
								const timestamp = nowTime + 1800 * 1000;
								const body = JSON.parse(resp.body);
								localStorage.setItem(
									weatherTag,
									JSON.stringify({
										timestamp: timestamp,
										data: body
									})
								);
								this.parseWeather(body);
							}
						},
						function() {}
					);
			} else {
				this.parseWeather(weather.data);
			}
		},
		parseWeather(w) {
			if (w.HeWeather6 && w.HeWeather6.length > 0 && w.HeWeather6[0].now) {
				const now = w.HeWeather6[0].now;
				const wa = [];
				wa.push(now.cond_txt);
				wa.push(now.tmp + "℃");
				wa.push(this.$t("relativeHumidity") + " " + now.hum + "%");
				wa.push(now.wind_dir + " " + now.wind_sc + this.$t("level"));
				this.weatherString = wa.join(" ");
			}
		},
		changeLang(lang) {
			const langTag =
				"signp-" + this.$store.state.startPointInfo.project_id + ":locale";
			const currentLang = localStorage.getItem(langTag);
			if (lang === currentLang) {
				return;
			}
			localStorage.setItem(langTag, lang);
			this.$i18n.locale = lang;
		},
		navPage(path) {
			this.$router.replace("/" + path);
		},
		bannerClick(item) {
			switch (item.bindType) {
				case "link":
					location.href = item.link;
					break;
				case "mapGid":
					this.navPage("details/" + item.mapGid);
					break;
				case "internal":
					this.navPage(item.link);
					break;
			}
		},
		loadPolygons() {
			this.markerList = [];
			const _this = this;
			const readStore = _this.$store.state.db
				.transaction("mapPolygons")
				.objectStore("mapPolygons")
				.getAll();
			readStore.onsuccess = function(e) {
				const r = e.target.result;
				if (r && r.length > 0) {
					for (let i = 0; i < r.length; i++) {
						r[i].lnglat = _this.changePoint(r[i].point);
						r[i].label = { content: _this.$t(r[i].name) };
						_this.markerList.push(r[i]);
					}
				}
			};
		}
	}
};
</script>
<style scoped>
.v-item-group.v-bottom-navigation .v-btn {
	min-width: 30px;
}

.toMapBtn {
	position: absolute;
	bottom: 70px;
	left: 10px;
	z-index: 3;
}
.amap {
	width: 100% !important;
}
.bg {
	width: 100%;
	height: 28em;
	top: 0;
	background: url("../assets/homeBg.jpg") no-repeat center center;
	background-size: cover;
}
.round {
	padding-top: 6px;
	border-radius: 50% 50% 0 0 !important;
	background-color: rgba(0, 0, 0, 0.5);
}
.roundBorder {
	border-radius: 50% 50% 0 0 !important;
	background-color: rgba(0, 0, 0, 0);
}
.roundBorder1 {
	margin-top: 10px;
	padding-top: 20px;
	border-radius: 50% 50% 0 0 !important;
	background-color: rgba(0, 0, 0, 0);
}
.roundBg {
	border-radius: 50% 50% 0 0 !important;
}
.v-list--three-line .v-list-item,
.v-list-item--three-line {
	min-height: 10px;
}
</style>
