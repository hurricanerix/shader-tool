// Copyright 2016 Richard Hawkins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package model manages reading 3D models.
package model

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

var Formats = [...]string{
	"format ascii 1.0",
}

type Model struct {
	Format      string
	VertexCount int
	FaceCount   int
	VertexData  []float32
	FaceData    []uint32
}

func New() Model {
	m := Model{}
	return m

}

func (m *Model) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	if err := m.readHeader(scanner); err != nil {
		return fmt.Errorf("could not read header:", err)
	}

	if err := m.readVertices(scanner); err != nil {
		return fmt.Errorf("could not read vertices:", err)
	}

	if err := m.readFaces(scanner); err != nil {
		return fmt.Errorf("could not read faces:", err)
	}

	return nil
}

func (m Model) String() string {
	msg := "PLY{\n"
	msg += fmt.Sprintf("  Format: %s\n", m.Format)
	msg += fmt.Sprintf("  VertexCount: %d\n", m.VertexCount)
	msg += fmt.Sprintf("  FaceCount: %d\n", m.FaceCount)
	msg += fmt.Sprintf("  VertexData: %v\n", m.VertexData)
	msg += fmt.Sprintf("  FaceData: %v\n", m.FaceData)
	msg += "}"
	return msg
}

func (m *Model) readHeader(scanner *bufio.Scanner) error {
	// Magic
	scanner.Scan()
	if magic := scanner.Text(); magic != "ply" {
		return fmt.Errorf("invalid ply file, expected 'ply' got '%s'", magic)
	}

	// format/version
	scanner.Scan()
	m.Format = scanner.Text()
	supported := false
	for i := range Formats {
		if m.Format == Formats[i] {
			supported = true
		}
	}
	if !supported {
		return fmt.Errorf("unsupported format: %s", m.Format)
	}

	var line string
	var t string
	var c int
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "end_header") {
			break
		}
		if strings.HasPrefix(line, "element") {
			if _, err := fmt.Sscan(line, &t, &t, &c); err != nil {
				return fmt.Errorf("trouble scanning header: %s", err)
			}
			if t == "vertex" {
				m.VertexCount = c
			} else if t == "face" {
				m.FaceCount = c
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (m *Model) readVertices(scanner *bufio.Scanner) error {
	m.VertexData = make([]float32, 8*m.VertexCount)
	var line string
	var x, y, z, nx, ny, nz, s, t float32
	for i := 0; i < m.VertexCount; i++ {
		scanner.Scan()
		line = scanner.Text()
		if _, err := fmt.Sscan(line, &x, &y, &z, &nx, &ny, &nz, &s, &t); err != nil {
			return fmt.Errorf("trouble scanning vertex: %s", err)
		}
		m.VertexData[(i*8)+0] = x
		m.VertexData[(i*8)+1] = y
		m.VertexData[(i*8)+2] = z
		m.VertexData[(i*8)+3] = nx
		m.VertexData[(i*8)+4] = ny
		m.VertexData[(i*8)+5] = nz
		m.VertexData[(i*8)+6] = s
		m.VertexData[(i*8)+7] = t
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("trouble scanning vertices: %s", err)
		}
	}
	return nil
}

func (m *Model) readFaces(scanner *bufio.Scanner) error {
	m.FaceData = make([]uint32, 3*m.FaceCount)
	var line string
	var c, v0, v1, v2 uint32
	for i := 0; i < m.FaceCount; i++ {
		scanner.Scan()
		line = scanner.Text()
		if _, err := fmt.Sscan(line, &c, &v0, &v1, &v2); err != nil {
			return fmt.Errorf("trouble scanning faces: %s", err)
		}
		m.FaceData[(i*3)+0] = v0
		m.FaceData[(i*3)+1] = v1
		m.FaceData[(i*3)+2] = v2
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("trouble scanning vertices: %s", err)
		}
	}
	return nil
}
