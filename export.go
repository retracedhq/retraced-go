package retraced

import (
	"context"
	"encoding/csv"
	"io"
)

// ExportCSV writes all events matching a query to w as CSV records
func (c *Client) ExportCSV(ctx context.Context, w io.Writer, sq *StructuredQuery, mask *EventNodeMask) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	events := make(chan *EventNode)
	errors := make(chan error, 1)
	out := csv.NewWriter(w)

	defer func() {
		out.Flush()
		cancel()
	}()

	if err := out.Write(mask.CSVHeaders()); err != nil {
		return err
	}

	go func() {
		stream, err := c.NewStream(sq, mask)
		if err != nil {
			errors <- err
			return
		}
		for {
			e, err := stream.Read()
			if err == io.EOF {
				close(events)
				return
			}
			if err != nil {
				errors <- err
				return
			}
			select {
			case events <- e:
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return nil
			}
			if err := out.Write(mask.CSVRow(e)); err != nil {
				return err
			}
		case err := <-errors:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}
