import React, {useState, useEffect} from 'react';
import ReactMarkdown from 'react-markdown';
import {isMobile, } from "react-device-detect";

import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Divider from '@material-ui/core/Divider';

// Component for the blog. Imported by ./App
const Blog = (props) => {
  const [blog, setBlog] = useState({});
  const [editing, setEditing] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [shortDescription, setShortDescription] = useState("");
  const [headerImage, setHeaderImage] = useState("");
  const [change, setChange] = useState("");
  const [writerName, setWriterName] = useState("");

	// A function that saves to the backend
	const onSave = (field, value) => {
		// Only update if new data
		if (value === blog[field]) {
			return
		}

		const timeNow = Date.now()/1000
		blog["updatedTime"] = timeNow
		if (blog.createdTime === undefined) {
			blog["createdTime"] = timeNow
		}

		blog[field] = value
		setBlog(blog)
	}

	// A component that updates your states only if they're actually updated
	const handleAdminSave = () => {
		var changed = false
		if (title.length > 0 && title !== blog.title) {
			onSave("title", title)
			changed = true
		}

		if (content.length > 0 && content !== blog.content) {
			onSave("content", content)
			changed = true
		}

		if (shortDescription.length > 0 && shortDescription !== blog.shortDescription) {
			onSave("shortDescription", shortDescription)
			changed = true
		}

		if (headerImage.length > 0 && headerImage !== blog.headerImage) {
			onSave("headerImage", headerImage)
			changed = true
		}

		if (writerName.length > 0 && writerName !== blog.writerName) {
			onSave("writerName", writerName)
			changed = true
		}

		// Workaround for realtime state update lol
		if (changed) {
			if (change.length > 10) {
				setChange("a")
			} else {
  			setChange(change+"a")
			}
		}
	}

	// Matches ctrl+s
	const keypress = (event) => {
		if (event.keyCode === 83 && event.ctrlKey) {
			event.preventDefault()
			const item = document.getElementById("saveButton")
			if (item !== null) {
				item.click()
			}
		}
	}
    		
	// Starts the keypress
	useEffect(() => {
		document.addEventListener("keydown", keypress, false)
	})

	// The actual editor to show
	const editingContent = editing ?
		<div style={{flex: 1, margin: isMobile ? 10 : 20, padding: 10, borderRight: "1px solid rgba(0,0,0,0.1)"}}>
			<Typography variant="h3">
				Edit	
			</Typography>
			<Divider style={{marginTop: 15, marginBottom: 15}}/>
			<TextField
				color="primary"
				required
				InputProps={{
					style:{
						height: "50px", 
						fontSize: "1em",
					},
				}}
				fullWidth={true}
				defaultValue={blog.title}
				label="Display title"
				margin="normal"
				variant="outlined"
				onChange={(event) => {
					setTitle(event.target.value)
				}}
			/>
			<TextField
				color="primary"
				required
				InputProps={{
					style:{
						height: "50px", 
						fontSize: "1em",
					},
				}}
				fullWidth={true}
				defaultValue={blog.shortDescription}
				label="Short description"
				margin="normal"
				variant="outlined"
				onChange={(event) => {
					setShortDescription(event.target.value)
				}}
			/>
			<TextField
				color="primary"
				required
				InputProps={{
					style:{
						height: "50px", 
						fontSize: "1em",
					},
				}}
				fullWidth={true}
				defaultValue={blog.headerImage}
				label="Header image"
				margin="normal"
				variant="outlined"
				onChange={(event) => {
					setHeaderImage(event.target.value)
				}}
			/>
			<TextField
				color="primary"
				required
				InputProps={{
					style:{
						height: "50px", 
						fontSize: "1em",
					},
				}}
				fullWidth={true}
				defaultValue={blog.writerName}
				label="Writer name"
				margin="normal"
				variant="outlined"
				onChange={(event) => {
					setWriterName(event.target.value)
				}}
			/>
			<Divider style={{marginTop: 20, marginBottom: 20}}/>
			<div style={{position: "sticky", top: 100}}>
				<Typography>Write using <a href="https://en.wikipedia.org/wiki/Markdown">markdown</a></Typography>
				<TextField
					color="primary"
					required
					InputProps={{
						style:{
							fontSize: "1em",
							height: "80%",
							backgroundColor: "white",
						},
					}}
					multiline
					fullWidth={true}
					defaultValue={blog.content}
					label="Content"
					margin="normal"
					variant="outlined"
					onChange={(event) => {
						setContent(event.target.value)
					}}
				/>
				<Button id="saveButton" color="secondary" style={{marginTop: 50, minWidth: 300, height: 50, borderRadius: 10, fontSize: 15}} variant="contained" onClick={() => {handleAdminSave()}}>Save (ctrl+s)</Button>
			</div>
		</div>
		: null

	// Style for the component
	const outerStyle = {
		display: "flex", 
		marginLeft: (isMobile ? 0 : "auto"), 
		marginRight: (isMobile ? 0 : "auto"), 
		width: (isMobile ? "100%" : editing ? "100%" : 850)
	}

	const htmlContent =
		<div style={outerStyle}>	
			{editingContent}
			<div style={{flex: 1, padding: isMobile ? 15 : 0,}}>
				<Button color="secondary" style={{marginTop: 50, marginBottom: 20, minWidth: 300, height: 50, borderRadius: 10, fontSize: 15}} variant="contained" onClick={() => setEditing(!editing)}>{editing ? "STOP EDIT" : "EDIT"}</Button>  
				{editing ? 
				<div>
					<Typography variant="h3">
						Blog	
					</Typography>
					<Divider style={{marginTop: 15, marginBottom: 15}}/>
				</div>
				: null }
				<img alt={blog.title} src={blog.headerImage} style={{width: "100%",}}/>
				<Typography style={{fontSize: 25}}>
					<ReactMarkdown 
						escapeHtml={false}
						source={blog.content} 
					/>
				</Typography>
			</div>
		</div> 

	return htmlContent
}

export default Blog
