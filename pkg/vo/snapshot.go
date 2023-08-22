package vo

import "time"

type Snashot struct {
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	Location               string `json:"location"`
	Compress               string `json:"compress"`
	MaxRestoreBytesPerSec  string `json:"max_restore_bytes_per_sec"`
	MaxSnapshotBytesPerSec string `json:"max_snapshot_bytes_per_sec"`
	ChunkSize              string `json:"chunk_size"`
	Readonly               string `json:"readonly"`
}

type SnapshotRepository struct {
	Type     string                     `json:"type"`
	Settings SnapshotRepositorySettings `json:"settings"`
}

type SnapshotRepositorySettings struct {
	Location               string `json:"location"`
	MaxRestoreBytesPerSec  string `json:"max_restore_bytes_per_sec"`
	Readonly               string `json:"readonly"`
	Compress               string `json:"compress"`
	MaxSnapshotBytesPerSec string `json:"max_snapshot_bytes_per_sec"`
}

type Snapshot struct {
	Id               string `json:"id"`
	Status           string `json:"status"`
	StartEpoch       string `json:"start_epoch"`
	StartTime        string `json:"start_time"`
	EndEpoch         string `json:"end_epoch"`
	EndTime          string `json:"end_time"`
	Duration         string `json:"duration"`
	Indices          string `json:"indices"`
	SuccessfulShards string `json:"successful_shards"`
	FailedShards     string `json:"failed_shards"`
	TotalShards      string `json:"total_shards"`
}

type SnapshotDetail struct {
	Snapshots []struct {
		Snapshot           string        `json:"snapshot"`
		Uuid               string        `json:"uuid"`
		VersionId          int           `json:"version_id"`
		Version            string        `json:"version"`
		Indices            []string      `json:"indices"`
		IncludeGlobalState bool          `json:"include_global_state"`
		State              string        `json:"state"`
		StartTime          time.Time     `json:"start_time"`
		StartTimeInMillis  int64         `json:"start_time_in_millis"`
		EndTime            time.Time     `json:"end_time"`
		EndTimeInMillis    int64         `json:"end_time_in_millis"`
		DurationInMillis   int           `json:"duration_in_millis"`
		Failures           []interface{} `json:"failures"`
		Shards             struct {
			Total      int `json:"total"`
			Failed     int `json:"failed"`
			Successful int `json:"successful"`
		} `json:"shards"`
	} `json:"snapshots"`
}

type SnapshotStatus struct {
	Snapshots []struct {
		Snapshot           string `json:"snapshot"`
		Repository         string `json:"repository"`
		Uuid               string `json:"uuid"`
		State              string `json:"state"`
		IncludeGlobalState bool   `json:"include_global_state"`
		ShardsStats        struct {
			Initializing int `json:"initializing"`
			Started      int `json:"started"`
			Finalizing   int `json:"finalizing"`
			Done         int `json:"done"`
			Failed       int `json:"failed"`
			Total        int `json:"total"`
		} `json:"shards_stats"`
		Stats struct {
			Incremental struct {
				FileCount   int `json:"file_count"`
				SizeInBytes int `json:"size_in_bytes"`
			} `json:"incremental"`
			Total struct {
				FileCount   int `json:"file_count"`
				SizeInBytes int `json:"size_in_bytes"`
			} `json:"total"`
			StartTimeInMillis    int64 `json:"start_time_in_millis"`
			TimeInMillis         int   `json:"time_in_millis"`
			NumberOfFiles        int   `json:"number_of_files"`
			ProcessedFiles       int   `json:"processed_files"`
			TotalSizeInBytes     int   `json:"total_size_in_bytes"`
			ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
		} `json:"stats"`
		Indices struct {
			Test2 struct {
				ShardsStats struct {
					Initializing int `json:"initializing"`
					Started      int `json:"started"`
					Finalizing   int `json:"finalizing"`
					Done         int `json:"done"`
					Failed       int `json:"failed"`
					Total        int `json:"total"`
				} `json:"shards_stats"`
				Stats struct {
					Incremental struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"incremental"`
					Total struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"total"`
					StartTimeInMillis    int64 `json:"start_time_in_millis"`
					TimeInMillis         int   `json:"time_in_millis"`
					NumberOfFiles        int   `json:"number_of_files"`
					ProcessedFiles       int   `json:"processed_files"`
					TotalSizeInBytes     int   `json:"total_size_in_bytes"`
					ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
				} `json:"stats"`
				Shards struct {
					Field1 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"0"`
				} `json:"shards"`
			} `json:"test2"`
			Tt2 struct {
				ShardsStats struct {
					Initializing int `json:"initializing"`
					Started      int `json:"started"`
					Finalizing   int `json:"finalizing"`
					Done         int `json:"done"`
					Failed       int `json:"failed"`
					Total        int `json:"total"`
				} `json:"shards_stats"`
				Stats struct {
					Incremental struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"incremental"`
					Total struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"total"`
					StartTimeInMillis    int64 `json:"start_time_in_millis"`
					TimeInMillis         int   `json:"time_in_millis"`
					NumberOfFiles        int   `json:"number_of_files"`
					ProcessedFiles       int   `json:"processed_files"`
					TotalSizeInBytes     int   `json:"total_size_in_bytes"`
					ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
				} `json:"stats"`
				Shards struct {
					Field1 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"0"`
				} `json:"shards"`
			} `json:"tt2"`
			Ttt struct {
				ShardsStats struct {
					Initializing int `json:"initializing"`
					Started      int `json:"started"`
					Finalizing   int `json:"finalizing"`
					Done         int `json:"done"`
					Failed       int `json:"failed"`
					Total        int `json:"total"`
				} `json:"shards_stats"`
				Stats struct {
					Incremental struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"incremental"`
					Total struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"total"`
					StartTimeInMillis    int64 `json:"start_time_in_millis"`
					TimeInMillis         int   `json:"time_in_millis"`
					NumberOfFiles        int   `json:"number_of_files"`
					ProcessedFiles       int   `json:"processed_files"`
					TotalSizeInBytes     int   `json:"total_size_in_bytes"`
					ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
				} `json:"stats"`
				Shards struct {
					Field1 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"0"`
				} `json:"shards"`
			} `json:"ttt"`
			Articles struct {
				ShardsStats struct {
					Initializing int `json:"initializing"`
					Started      int `json:"started"`
					Finalizing   int `json:"finalizing"`
					Done         int `json:"done"`
					Failed       int `json:"failed"`
					Total        int `json:"total"`
				} `json:"shards_stats"`
				Stats struct {
					Incremental struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"incremental"`
					Total struct {
						FileCount   int `json:"file_count"`
						SizeInBytes int `json:"size_in_bytes"`
					} `json:"total"`
					StartTimeInMillis    int64 `json:"start_time_in_millis"`
					TimeInMillis         int   `json:"time_in_millis"`
					NumberOfFiles        int   `json:"number_of_files"`
					ProcessedFiles       int   `json:"processed_files"`
					TotalSizeInBytes     int   `json:"total_size_in_bytes"`
					ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
				} `json:"stats"`
				Shards struct {
					Field1 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"0"`
					Field2 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"1"`
					Field3 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"2"`
					Field4 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"3"`
					Field5 struct {
						Stage string `json:"stage"`
						Stats struct {
							Incremental struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"incremental"`
							Total struct {
								FileCount   int `json:"file_count"`
								SizeInBytes int `json:"size_in_bytes"`
							} `json:"total"`
							StartTimeInMillis    int64 `json:"start_time_in_millis"`
							TimeInMillis         int   `json:"time_in_millis"`
							NumberOfFiles        int   `json:"number_of_files"`
							ProcessedFiles       int   `json:"processed_files"`
							TotalSizeInBytes     int   `json:"total_size_in_bytes"`
							ProcessedSizeInBytes int   `json:"processed_size_in_bytes"`
						} `json:"stats"`
					} `json:"4"`
				} `json:"shards"`
			} `json:"articles"`
		} `json:"indices"`
	} `json:"snapshots"`
}
