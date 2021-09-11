package main

type Output map[string]Environment

type Environment map[string]Region

type Region map[string]Server

type Server struct {
	Images map[string]string `json:"images,omitempty"`
	Tags   map[string]string `json:"tags,omitempty"`
}
