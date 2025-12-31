

data_dir := Path.home() / '.poopy'
data_file := data_dir / 'expenses.json'

func init_data_file() {
	data_dir.mkdir(exist_ok=True)
	if !data_file.exists() {

	}
}