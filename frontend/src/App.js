import React, {useState } from 'react';

const App = (props) => {
  const name = "@frikkylikeme"
  const url = "https://medium.com"
 
  return (
    <div>
      <Functionname url={url} name={name} {...props} /> 
    </div>
  );
};

const Functionname = props => {
  const {url, name} = props
	const [file, setFile] = useState("");

	const uploadFile = () => {
		fetch(url, {
			method: "POST",
			data: file,
		})
		.then(response => {
			console.log(response)
		})
		.catch(error => {
			console.log(error)
		})
	}

	return (
		<div>
			<p>{name}</p>
			<input type="file" name="fieldname" onChange={event => setFile(event.target.files[0])} />
			<button onClick={file => uploadFile(file)}>Upload</button>
		</div>
	) 
};

export default App;