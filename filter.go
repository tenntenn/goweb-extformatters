/*
A plugin of goweb which provides extra formatters and decoders.
https://github.com/tenntenn/goweb-extformatters

Copyright (c) 2012, Takuya Ueda. 
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice,
  this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.
* Neither the name of the author nor the names of its contributors may be used
  to endorse or promote products derived from this software
  without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package extformatters

import (
    "code.google.com/p/goweb/goweb"
    "io/ioutil"
    "strings"
)

type Filter interface {
    Filter([]uint8 input) ([]uint8, error)
}

type FilterFunc func([]uint8 input) ([]uint8,error)
func (f FilterFunc) Filter([]uint8 input) ([]uint8, error) {
   return f(input)
}

type FilteredDecoder struct {
    decoder goweb.Decoder
    filters []Filter
}

func NewFilteredDecoder(decoder goweb.Decoder, filters ...Filter) *FilteredDecoder {
    cpy := make([]Filter, len(filters))
    copy(cpy, filters)

    return &FilteredDecoder{decoder, cpy}
}

func (d *FilteredDecoder) Unmarshal(cx *goweb.Context, v interface{}) (err error) {
    body := ioutil.ReadAll(cx.Request.Body)
    for _, f := range d.filters {
        if f == nil {
            continue
        }

        body, err = f.Filter(body)
        if err != nil {
            return nil
        }
    }

    return d.decoder.Unmarshal(body, v)
}
